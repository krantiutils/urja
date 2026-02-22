package nutrition

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"strings"
	"time"
)

// Service handles nutrition business logic.
type Service struct {
	repo   *Repository
	logger *slog.Logger
}

// NewService creates a new nutrition service.
func NewService(repo *Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// --- Food Items ---

// SearchFoods searches the food catalog with optional category filter.
func (s *Service) SearchFoods(ctx context.Context, userID, orgID, query, category string, limit int) ([]FoodItem, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	query = strings.TrimSpace(query)
	if query == "" {
		return s.repo.ListFoodItems(ctx, orgID, userID, category, limit)
	}
	return s.repo.SearchFoodItems(ctx, orgID, userID, query, category, limit)
}

// GetFoodItem retrieves a single food item.
func (s *Service) GetFoodItem(ctx context.Context, id string) (*FoodItem, error) {
	return s.repo.GetFoodItem(ctx, id)
}

// --- Custom Foods ---

// CreateCustomFoodInput holds input for creating a user-defined food item.
type CreateCustomFoodInput struct {
	Name            string  `json:"name"`
	NameNe          string  `json:"name_ne,omitempty"`
	Category        string  `json:"category,omitempty"`
	CaloriesPer100g float64 `json:"calories_per_100g"`
	ProteinPer100g  float64 `json:"protein_per_100g"`
	CarbsPer100g    float64 `json:"carbs_per_100g"`
	FatPer100g      float64 `json:"fat_per_100g"`
	FiberPer100g    float64 `json:"fiber_per_100g,omitempty"`
	ServingSizeG    float64 `json:"serving_size_g,omitempty"`
	ServingLabel    string  `json:"serving_label,omitempty"`
	Barcode         string  `json:"barcode,omitempty"`
}

// CreateCustomFood creates a user-defined food item.
func (s *Service) CreateCustomFood(ctx context.Context, userID string, input *CreateCustomFoodInput) (*FoodItem, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	input.Name = name
	if input.CaloriesPer100g < 0 {
		return nil, fmt.Errorf("calories_per_100g must be >= 0")
	}
	if input.ProteinPer100g < 0 {
		return nil, fmt.Errorf("protein_per_100g must be >= 0")
	}
	if input.CarbsPer100g < 0 {
		return nil, fmt.Errorf("carbs_per_100g must be >= 0")
	}
	if input.FatPer100g < 0 {
		return nil, fmt.Errorf("fat_per_100g must be >= 0")
	}

	f, err := s.repo.CreateCustomFood(ctx, userID, input)
	if err != nil {
		return nil, err
	}

	s.logger.Info("custom food created", "food_id", f.ID, "user_id", userID, "name", input.Name)
	return f, nil
}

// LookupBarcode looks up a food item by barcode. Checks local DB first, then Open Food Facts.
func (s *Service) LookupBarcode(ctx context.Context, barcode string) (*FoodItem, error) {
	barcode = strings.TrimSpace(barcode)
	if barcode == "" {
		return nil, fmt.Errorf("barcode is required")
	}

	// Check local DB first.
	f, err := s.repo.GetFoodByBarcode(ctx, barcode)
	if err == nil {
		return f, nil
	}

	// Not found locally, try Open Food Facts API.
	offProduct, err := fetchOpenFoodFacts(barcode)
	if err != nil {
		return nil, fmt.Errorf("food not found for barcode %s", barcode)
	}

	// Save to local DB as a verified food item (system-created from OFF).
	saved, err := s.repo.CreateVerifiedFoodFromBarcode(ctx, offProduct.Name,
		offProduct.CaloriesPer100g, offProduct.ProteinPer100g, offProduct.CarbsPer100g,
		offProduct.FatPer100g, offProduct.FiberPer100g, barcode)
	if err != nil {
		return nil, fmt.Errorf("saving barcode food: %w", err)
	}

	s.logger.Info("food imported from Open Food Facts", "food_id", saved.ID, "barcode", barcode, "name", saved.Name)
	return saved, nil
}

// GetNutritionStreak retrieves the nutrition logging streak for a user.
func (s *Service) GetNutritionStreak(ctx context.Context, userID string) (*NutritionStreak, error) {
	return s.repo.GetNutritionStreak(ctx, userID)
}

// --- Food Logs ---

var validMealTypes = map[string]bool{
	"breakfast": true,
	"lunch":     true,
	"dinner":    true,
	"snack":     true,
}

// CreateFoodLogInput holds input for logging food.
type CreateFoodLogInput struct {
	OrgID         string  `json:"organization_id"`
	FoodItemID    string  `json:"food_item_id"`
	MealType      string  `json:"meal_type"`
	QuantityGrams float64 `json:"quantity_grams"`
	LoggedDate    string  `json:"logged_date"`
	Notes         string  `json:"notes"`
}

// CreateFoodLog logs a food item, pre-calculating macros from the food item's per-100g values.
func (s *Service) CreateFoodLog(ctx context.Context, userID string, input *CreateFoodLogInput) (*FoodLog, error) {
	if input.FoodItemID == "" {
		return nil, fmt.Errorf("food_item_id is required")
	}
	if !validMealTypes[input.MealType] {
		return nil, fmt.Errorf("meal_type must be breakfast, lunch, dinner, or snack")
	}
	if input.QuantityGrams <= 0 {
		return nil, fmt.Errorf("quantity_grams must be positive")
	}
	if input.LoggedDate == "" {
		return nil, fmt.Errorf("logged_date is required (YYYY-MM-DD)")
	}
	if _, err := time.Parse("2006-01-02", input.LoggedDate); err != nil {
		return nil, fmt.Errorf("logged_date must be in YYYY-MM-DD format")
	}

	// Lookup food item to get per-100g values.
	food, err := s.repo.GetFoodItem(ctx, input.FoodItemID)
	if err != nil {
		return nil, fmt.Errorf("food item not found")
	}

	// Pre-calculate macros.
	calories := food.CaloriesPer100g * input.QuantityGrams / 100
	protein := food.ProteinPer100g * input.QuantityGrams / 100
	carbs := food.CarbsPer100g * input.QuantityGrams / 100
	fat := food.FatPer100g * input.QuantityGrams / 100

	l, err := s.repo.CreateFoodLog(ctx, userID, input.OrgID, input.FoodItemID, input.MealType,
		input.QuantityGrams, calories, protein, carbs, fat,
		input.LoggedDate, strings.TrimSpace(input.Notes))
	if err != nil {
		return nil, err
	}

	l.FoodItem = food

	// Update nutrition streak.
	if _, err := s.repo.UpdateNutritionStreak(ctx, userID, input.LoggedDate); err != nil {
		s.logger.Warn("failed to update nutrition streak", "error", err, "user_id", userID)
	}

	s.logger.Info("food logged", "log_id", l.ID, "user_id", userID, "food_item_id", input.FoodItemID, "org_id", input.OrgID)
	return l, nil
}

// ListFoodLogs retrieves food logs for a user on a given date.
func (s *Service) ListFoodLogs(ctx context.Context, userID, orgID, date, mealType string) ([]FoodLog, error) {
	if date == "" {
		return nil, fmt.Errorf("date is required (YYYY-MM-DD)")
	}
	return s.repo.ListFoodLogs(ctx, userID, orgID, date, mealType)
}

// DeleteFoodLog deletes a food log entry owned by the user.
func (s *Service) DeleteFoodLog(ctx context.Context, userID, logID string) error {
	if err := s.repo.DeleteFoodLog(ctx, userID, logID); err != nil {
		return err
	}
	s.logger.Info("food log deleted", "log_id", logID, "user_id", userID)
	return nil
}

// --- Daily / Weekly Summary ---

// GetDailySummary returns aggregated nutrition data for a single day, grouped by meal.
func (s *Service) GetDailySummary(ctx context.Context, userID, orgID, date string) (*DailySummary, error) {
	if date == "" {
		return nil, fmt.Errorf("date is required (YYYY-MM-DD)")
	}

	logs, err := s.repo.GetDailySummaryLogs(ctx, userID, orgID, date)
	if err != nil {
		return nil, err
	}

	// Group by meal type and compute totals.
	mealMap := make(map[string]*MealSummary)
	mealOrder := []string{"breakfast", "lunch", "dinner", "snack"}

	summary := &DailySummary{Date: date}

	for _, l := range logs {
		summary.TotalCalories += l.Calories
		summary.TotalProtein += l.Protein
		summary.TotalCarbs += l.Carbs
		summary.TotalFat += l.Fat

		ms, ok := mealMap[l.MealType]
		if !ok {
			ms = &MealSummary{MealType: l.MealType}
			mealMap[l.MealType] = ms
		}
		ms.Calories += l.Calories
		ms.Protein += l.Protein
		ms.Carbs += l.Carbs
		ms.Fat += l.Fat
		ms.Items = append(ms.Items, l)
	}

	// Build meals in standard order.
	for _, mt := range mealOrder {
		if ms, ok := mealMap[mt]; ok {
			summary.Meals = append(summary.Meals, *ms)
		}
	}

	return summary, nil
}

// GetWeeklySummary returns daily aggregated totals for 7 days starting from a date.
func (s *Service) GetWeeklySummary(ctx context.Context, userID, orgID, fromDate string) ([]DailySummary, error) {
	if fromDate == "" {
		return nil, fmt.Errorf("from date is required (YYYY-MM-DD)")
	}
	return s.repo.GetWeeklySummary(ctx, userID, orgID, fromDate)
}

// --- Daily Dashboard ---

// GetDailyDashboard returns a combined daily nutrition dashboard with calories, macros, water, meals, and streak.
func (s *Service) GetDailyDashboard(ctx context.Context, userID, orgID, date string) (*DailyDashboard, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	// 1. Get daily summary (meals and food totals).
	dailySummary, err := s.GetDailySummary(ctx, userID, orgID, date)
	if err != nil {
		return nil, fmt.Errorf("getting daily summary: %w", err)
	}

	// 2. Get nutrition goal (use defaults if not found).
	var calorieGoal, proteinGoal, carbsGoal, fatGoal float64
	goal, err := s.GetNutritionGoal(ctx, userID, orgID)
	if err != nil {
		// Use defaults when goal not found.
		s.logger.Debug("nutrition goal not found, using defaults", "user_id", userID, "error", err)
		calorieGoal = 2000
		proteinGoal = 150
		carbsGoal = 200
		fatGoal = 67
	} else {
		calorieGoal = goal.CalorieGoal
		proteinGoal = goal.ProteinGoalG
		carbsGoal = goal.CarbsGoalG
		fatGoal = goal.FatGoalG
	}

	// 3. Get nutrition streak (use 0/0 if not found).
	var currentStreak, longestStreak int
	streak, err := s.repo.GetNutritionStreak(ctx, userID)
	if err != nil {
		s.logger.Debug("nutrition streak not found, using defaults", "user_id", userID, "error", err)
	} else {
		currentStreak = streak.CurrentStreak
		longestStreak = streak.LongestStreak
	}

	// 4. Get water total for the date.
	waterTotal, err := s.repo.GetWaterTotal(ctx, userID, date)
	if err != nil {
		s.logger.Debug("water total not found, using 0", "user_id", userID, "error", err)
		waterTotal = 0
	}

	// 5. Get user water goal.
	waterGoal, err := s.repo.GetUserWaterGoal(ctx, userID)
	if err != nil {
		s.logger.Debug("user water goal not found, using default", "user_id", userID, "error", err)
		waterGoal = 2500
	}

	// 6. Combine into dashboard response.
	remaining := calorieGoal - dailySummary.TotalCalories
	if remaining < 0 {
		remaining = 0
	}

	dashboard := &DailyDashboard{
		Date: date,
		Calories: DashboardCalories{
			Consumed:  dailySummary.TotalCalories,
			Goal:      calorieGoal,
			Remaining: remaining,
		},
		Macros: DashboardMacros{
			Protein: DashboardMacro{G: dailySummary.TotalProtein, Goal: proteinGoal},
			Carbs:   DashboardMacro{G: dailySummary.TotalCarbs, Goal: carbsGoal},
			Fat:     DashboardMacro{G: dailySummary.TotalFat, Goal: fatGoal},
		},
		Water: DashboardWater{
			ML:   waterTotal,
			Goal: waterGoal,
		},
		Meals: dailySummary.Meals,
		Streak: DashboardStreak{
			Current: currentStreak,
			Longest: longestStreak,
		},
	}

	// Ensure meals is never nil in JSON output.
	if dashboard.Meals == nil {
		dashboard.Meals = []MealSummary{}
	}

	return dashboard, nil
}

// --- Nutrition Goals (TDEE calculation) ---

var activityMultipliers = map[string]float64{
	"sedentary":   1.2,
	"light":       1.375,
	"moderate":    1.55,
	"active":      1.725,
	"very_active": 1.9,
}

var goalAdjustments = map[string]float64{
	"lose_weight":  -500,
	"maintain":     0,
	"build_muscle": 300,
}

// SetNutritionGoalInput holds input for setting a TDEE-based nutrition goal.
type SetNutritionGoalInput struct {
	OrgID         string  `json:"organization_id"`
	WeightKg      float64 `json:"weight_kg"`
	HeightCm      float64 `json:"height_cm"`
	Age           int     `json:"age"`
	Gender        string  `json:"gender"`
	ActivityLevel string  `json:"activity_level"`
	GoalType      string  `json:"goal_type"`
}

// SetNutritionGoal calculates TDEE and sets the nutrition goal for a user.
func (s *Service) SetNutritionGoal(ctx context.Context, userID string, input *SetNutritionGoalInput) (*NutritionGoal, error) {
	if input.WeightKg <= 0 {
		return nil, fmt.Errorf("weight_kg must be positive")
	}
	if input.HeightCm <= 0 {
		return nil, fmt.Errorf("height_cm must be positive")
	}
	if input.Age <= 0 {
		return nil, fmt.Errorf("age must be positive")
	}
	if input.Gender != "male" && input.Gender != "female" {
		return nil, fmt.Errorf("gender must be male or female")
	}

	multiplier, ok := activityMultipliers[input.ActivityLevel]
	if !ok {
		return nil, fmt.Errorf("activity_level must be sedentary, light, moderate, active, or very_active")
	}

	adjustment, ok := goalAdjustments[input.GoalType]
	if !ok {
		return nil, fmt.Errorf("goal_type must be lose_weight, maintain, or build_muscle")
	}

	// Mifflin-St Jeor equation.
	var offset float64
	if input.Gender == "male" {
		offset = 5
	} else {
		offset = -161
	}
	bmr := 10*input.WeightKg + 6.25*input.HeightCm - 5*float64(input.Age) + offset
	tdee := bmr * multiplier
	calorieGoal := math.Round(tdee + adjustment)

	// Default macros: Protein 30%, Carbs 40%, Fat 30%.
	proteinGoalG := math.Round(calorieGoal * 0.30 / 4)
	carbsGoalG := math.Round(calorieGoal * 0.40 / 4)
	fatGoalG := math.Round(calorieGoal * 0.30 / 9)

	g, err := s.repo.UpsertNutritionGoal(ctx, userID, input.OrgID,
		calorieGoal, proteinGoalG, carbsGoalG, fatGoalG,
		input.WeightKg, input.HeightCm, input.Age,
		input.Gender, input.ActivityLevel, input.GoalType)
	if err != nil {
		return nil, err
	}

	s.logger.Info("nutrition goal set", "user_id", userID, "org_id", input.OrgID,
		"calorie_goal", calorieGoal, "tdee", tdee, "bmr", bmr)
	return g, nil
}

// GetNutritionGoal retrieves the current nutrition goal for a user.
func (s *Service) GetNutritionGoal(ctx context.Context, userID, orgID string) (*NutritionGoal, error) {
	return s.repo.GetNutritionGoal(ctx, userID, orgID)
}

// --- Meal Templates ---

// CreateMealTemplateInput holds input for creating a meal template.
type CreateMealTemplateInput struct {
	OrgID    string          `json:"organization_id"`
	Name     string          `json:"name"`
	NameNe   string          `json:"name_ne"`
	MealType string          `json:"meal_type"`
	Items    json.RawMessage `json:"items"`
}

// CreateMealTemplate creates a saved meal template.
func (s *Service) CreateMealTemplate(ctx context.Context, userID string, input *CreateMealTemplateInput) (*MealTemplate, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if input.MealType != "" && !validMealTypes[input.MealType] {
		return nil, fmt.Errorf("meal_type must be breakfast, lunch, dinner, or snack")
	}

	items := input.Items
	if len(items) == 0 {
		items = json.RawMessage(`[]`)
	}

	t, err := s.repo.CreateMealTemplate(ctx, userID, input.OrgID, name, strings.TrimSpace(input.NameNe), input.MealType, items)
	if err != nil {
		return nil, err
	}

	s.logger.Info("meal template created", "template_id", t.ID, "user_id", userID, "org_id", input.OrgID)
	return t, nil
}

// ListMealTemplates retrieves all meal templates for a user in an org.
func (s *Service) ListMealTemplates(ctx context.Context, userID, orgID string) ([]MealTemplate, error) {
	return s.repo.ListMealTemplates(ctx, userID, orgID)
}

// DeleteMealTemplate deletes a meal template owned by the user.
func (s *Service) DeleteMealTemplate(ctx context.Context, userID, templateID string) error {
	if err := s.repo.DeleteMealTemplate(ctx, userID, templateID); err != nil {
		return err
	}
	s.logger.Info("meal template deleted", "template_id", templateID, "user_id", userID)
	return nil
}

// templateItem represents an item within a meal template's JSONB items array.
type templateItem struct {
	FoodItemID    string  `json:"food_item_id"`
	QuantityGrams float64 `json:"quantity_grams"`
}

// LogMealTemplateInput holds input for logging a meal template.
type LogMealTemplateInput struct {
	OrgID      string `json:"organization_id"`
	LoggedDate string `json:"logged_date"`
}

// LogMealTemplate logs all items from a meal template as individual food logs.
func (s *Service) LogMealTemplate(ctx context.Context, userID, templateID string, input *LogMealTemplateInput) ([]FoodLog, error) {
	if input.LoggedDate == "" {
		return nil, fmt.Errorf("logged_date is required (YYYY-MM-DD)")
	}

	// Fetch the template.
	tmpl, err := s.repo.GetMealTemplate(ctx, templateID)
	if err != nil {
		return nil, fmt.Errorf("meal template not found")
	}

	// Verify ownership.
	if tmpl.UserID != userID {
		return nil, fmt.Errorf("meal template not found")
	}

	// Parse the items JSONB.
	var items []templateItem
	if err := json.Unmarshal(tmpl.Items, &items); err != nil {
		return nil, fmt.Errorf("invalid meal template items: %w", err)
	}

	mealType := tmpl.MealType
	if mealType == "" {
		mealType = "snack"
	}

	var logs []FoodLog
	for _, item := range items {
		// Lookup food item.
		food, err := s.repo.GetFoodItem(ctx, item.FoodItemID)
		if err != nil {
			return nil, fmt.Errorf("food item %s not found", item.FoodItemID)
		}

		// Calculate macros.
		calories := food.CaloriesPer100g * item.QuantityGrams / 100
		protein := food.ProteinPer100g * item.QuantityGrams / 100
		carbs := food.CarbsPer100g * item.QuantityGrams / 100
		fat := food.FatPer100g * item.QuantityGrams / 100

		l, err := s.repo.CreateFoodLog(ctx, userID, input.OrgID, item.FoodItemID, mealType,
			item.QuantityGrams, calories, protein, carbs, fat, input.LoggedDate, "")
		if err != nil {
			return nil, fmt.Errorf("logging template item: %w", err)
		}

		l.FoodItem = food
		logs = append(logs, *l)
	}

	s.logger.Info("meal template logged", "template_id", templateID, "user_id", userID,
		"items_count", len(logs), "org_id", input.OrgID)
	return logs, nil
}
