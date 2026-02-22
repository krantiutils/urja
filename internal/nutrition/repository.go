package nutrition

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// FoodItem represents an item in the food catalog.
type FoodItem struct {
	ID              string    `json:"id"`
	OrgID           *string   `json:"organization_id"`
	Name            string    `json:"name"`
	NameNe          string    `json:"name_ne,omitempty"`
	Category        string    `json:"category,omitempty"`
	CaloriesPer100g float64   `json:"calories_per_100g"`
	ProteinPer100g  float64   `json:"protein_per_100g"`
	CarbsPer100g    float64   `json:"carbs_per_100g"`
	FatPer100g      float64   `json:"fat_per_100g"`
	FiberPer100g    float64   `json:"fiber_per_100g"`
	ServingSizeG    float64   `json:"serving_size_g,omitempty"`
	ServingLabel    string    `json:"serving_label,omitempty"`
	ServingLabelNe  string    `json:"serving_label_ne,omitempty"`
	IsVerified      bool      `json:"is_verified"`
	Barcode         string    `json:"barcode,omitempty"`
	CreatedBy       *string   `json:"created_by,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// FoodLog represents a single food log entry.
type FoodLog struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	OrgID         *string   `json:"organization_id"`
	FoodItemID    string    `json:"food_item_id"`
	FoodItem      *FoodItem `json:"food_item,omitempty"`
	MealType      string    `json:"meal_type"`
	QuantityGrams float64   `json:"quantity_grams"`
	Calories      float64   `json:"calories"`
	Protein       float64   `json:"protein"`
	Carbs         float64   `json:"carbs"`
	Fat           float64   `json:"fat"`
	LoggedDate    string    `json:"logged_date"`
	LoggedAt      time.Time `json:"logged_at"`
	Notes         string    `json:"notes,omitempty"`
}

// MealTemplate represents a saved meal combination.
type MealTemplate struct {
	ID        string          `json:"id"`
	UserID    string          `json:"user_id"`
	OrgID     *string         `json:"organization_id"`
	Name      string          `json:"name"`
	NameNe    string          `json:"name_ne,omitempty"`
	MealType  string          `json:"meal_type"`
	Items     json.RawMessage `json:"items"`
	CreatedAt time.Time       `json:"created_at"`
}

// NutritionGoal represents daily nutrition targets for a user.
type NutritionGoal struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	OrgID         *string   `json:"organization_id"`
	CalorieGoal   float64   `json:"calorie_goal"`
	ProteinGoalG  float64   `json:"protein_goal_g"`
	CarbsGoalG    float64   `json:"carbs_goal_g"`
	FatGoalG      float64   `json:"fat_goal_g"`
	WeightKg      float64   `json:"weight_kg"`
	HeightCm      float64   `json:"height_cm"`
	Age           int       `json:"age"`
	Gender        string    `json:"gender"`
	ActivityLevel string    `json:"activity_level"`
	GoalType      string    `json:"goal_type"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// DailySummary represents aggregated nutrition data for a single day.
type DailySummary struct {
	Date          string        `json:"date"`
	TotalCalories float64       `json:"total_calories"`
	TotalProtein  float64       `json:"total_protein"`
	TotalCarbs    float64       `json:"total_carbs"`
	TotalFat      float64       `json:"total_fat"`
	Meals         []MealSummary `json:"meals"`
}

// MealSummary represents aggregated nutrition data for a single meal type.
type MealSummary struct {
	MealType string    `json:"meal_type"`
	Calories float64   `json:"calories"`
	Protein  float64   `json:"protein"`
	Carbs    float64   `json:"carbs"`
	Fat      float64   `json:"fat"`
	Items    []FoodLog `json:"items"`
}

// NutritionStreak represents a user's food logging streak.
type NutritionStreak struct {
	UserID        string    `json:"user_id"`
	CurrentStreak int       `json:"current_streak"`
	LongestStreak int       `json:"longest_streak"`
	LastLogDate   *string   `json:"last_log_date,omitempty"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// DailyDashboard represents the combined daily nutrition dashboard response.
type DailyDashboard struct {
	Date     string            `json:"date"`
	Calories DashboardCalories `json:"calories"`
	Macros   DashboardMacros   `json:"macros"`
	Water    DashboardWater    `json:"water"`
	Meals    []MealSummary     `json:"meals"`
	Streak   DashboardStreak   `json:"streak"`
}

// DashboardCalories holds calorie tracking data for the dashboard.
type DashboardCalories struct {
	Consumed  float64 `json:"consumed"`
	Goal      float64 `json:"goal"`
	Remaining float64 `json:"remaining"`
}

// DashboardMacros holds macro tracking data for the dashboard.
type DashboardMacros struct {
	Protein DashboardMacro `json:"protein"`
	Carbs   DashboardMacro `json:"carbs"`
	Fat     DashboardMacro `json:"fat"`
}

// DashboardMacro holds a single macro nutrient tracking data.
type DashboardMacro struct {
	G    float64 `json:"g"`
	Goal float64 `json:"goal"`
}

// DashboardWater holds water intake tracking data for the dashboard.
type DashboardWater struct {
	ML   int `json:"ml"`
	Goal int `json:"goal"`
}

// DashboardStreak holds nutrition logging streak data for the dashboard.
type DashboardStreak struct {
	Current int `json:"current"`
	Longest int `json:"longest"`
}

// Repository handles nutrition data persistence.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new nutrition repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// --- Food Items ---

// ListFoodItems returns food items without a search query, ordered by name.
func (r *Repository) ListFoodItems(ctx context.Context, orgID, userID, category string, limit int) ([]FoodItem, error) {
	var sql string
	var args []interface{}
	var argIdx int

	if orgID == "" {
		sql = `SELECT id, organization_id, name, COALESCE(name_ne, ''), COALESCE(category, ''),
		               COALESCE(calories_per_100g, 0), COALESCE(protein_per_100g, 0),
		               COALESCE(carbs_per_100g, 0), COALESCE(fat_per_100g, 0),
		               COALESCE(fiber_per_100g, 0), COALESCE(serving_size_g, 0),
		               COALESCE(serving_label, ''), COALESCE(serving_label_ne, ''),
		               is_verified, COALESCE(barcode, ''), created_by, created_at, updated_at
		        FROM food_items
		        WHERE (organization_id IS NULL OR created_by = $1)`
		args = []interface{}{userID}
		argIdx = 2
	} else {
		sql = `SELECT id, organization_id, name, COALESCE(name_ne, ''), COALESCE(category, ''),
		               COALESCE(calories_per_100g, 0), COALESCE(protein_per_100g, 0),
		               COALESCE(carbs_per_100g, 0), COALESCE(fat_per_100g, 0),
		               COALESCE(fiber_per_100g, 0), COALESCE(serving_size_g, 0),
		               COALESCE(serving_label, ''), COALESCE(serving_label_ne, ''),
		               is_verified, COALESCE(barcode, ''), created_by, created_at, updated_at
		        FROM food_items
		        WHERE (organization_id = $1 OR organization_id IS NULL OR created_by = $2)`
		args = []interface{}{orgID, userID}
		argIdx = 3
	}

	if category != "" {
		sql += fmt.Sprintf(" AND category = $%d", argIdx)
		args = append(args, category)
		argIdx++
	}

	sql += fmt.Sprintf(" ORDER BY name ASC LIMIT $%d", argIdx)
	args = append(args, limit)

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("listing food items: %w", err)
	}
	defer rows.Close()

	var items []FoodItem
	for rows.Next() {
		var f FoodItem
		if err := rows.Scan(&f.ID, &f.OrgID, &f.Name, &f.NameNe, &f.Category,
			&f.CaloriesPer100g, &f.ProteinPer100g, &f.CarbsPer100g, &f.FatPer100g,
			&f.FiberPer100g, &f.ServingSizeG, &f.ServingLabel, &f.ServingLabelNe,
			&f.IsVerified, &f.Barcode, &f.CreatedBy, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning food item: %w", err)
		}
		items = append(items, f)
	}
	return items, nil
}

// SearchFoodItems searches the food catalog using trigram similarity.
func (r *Repository) SearchFoodItems(ctx context.Context, orgID, userID, query, category string, limit int) ([]FoodItem, error) {
	var sql string
	var args []interface{}
	var argIdx int
	var queryArgIdx int

	if orgID == "" {
		sql = `SELECT id, organization_id, name, COALESCE(name_ne, ''), COALESCE(category, ''),
		               COALESCE(calories_per_100g, 0), COALESCE(protein_per_100g, 0),
		               COALESCE(carbs_per_100g, 0), COALESCE(fat_per_100g, 0),
		               COALESCE(fiber_per_100g, 0), COALESCE(serving_size_g, 0),
		               COALESCE(serving_label, ''), COALESCE(serving_label_ne, ''),
		               is_verified, COALESCE(barcode, ''), created_by, created_at, updated_at
		        FROM food_items
		        WHERE (organization_id IS NULL OR created_by = $1)
		          AND (name ILIKE '%' || $2 || '%' OR name_ne ILIKE '%' || $2 || '%')`
		args = []interface{}{userID, query}
		argIdx = 3
		queryArgIdx = 2
	} else {
		sql = `SELECT id, organization_id, name, COALESCE(name_ne, ''), COALESCE(category, ''),
		               COALESCE(calories_per_100g, 0), COALESCE(protein_per_100g, 0),
		               COALESCE(carbs_per_100g, 0), COALESCE(fat_per_100g, 0),
		               COALESCE(fiber_per_100g, 0), COALESCE(serving_size_g, 0),
		               COALESCE(serving_label, ''), COALESCE(serving_label_ne, ''),
		               is_verified, COALESCE(barcode, ''), created_by, created_at, updated_at
		        FROM food_items
		        WHERE (organization_id = $1 OR organization_id IS NULL OR created_by = $2)
		          AND (name ILIKE '%' || $3 || '%' OR name_ne ILIKE '%' || $3 || '%')`
		args = []interface{}{orgID, userID, query}
		argIdx = 4
		queryArgIdx = 3
	}

	if category != "" {
		sql += fmt.Sprintf(" AND category = $%d", argIdx)
		args = append(args, category)
		argIdx++
	}

	sql += fmt.Sprintf(" ORDER BY similarity(name, $%d) DESC, name ASC LIMIT $%d", queryArgIdx, argIdx)
	args = append(args, limit)

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("searching food items: %w", err)
	}
	defer rows.Close()

	var items []FoodItem
	for rows.Next() {
		var f FoodItem
		if err := rows.Scan(&f.ID, &f.OrgID, &f.Name, &f.NameNe, &f.Category,
			&f.CaloriesPer100g, &f.ProteinPer100g, &f.CarbsPer100g, &f.FatPer100g,
			&f.FiberPer100g, &f.ServingSizeG, &f.ServingLabel, &f.ServingLabelNe,
			&f.IsVerified, &f.Barcode, &f.CreatedBy, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning food item: %w", err)
		}
		items = append(items, f)
	}
	return items, rows.Err()
}

// GetFoodItem retrieves a single food item by ID.
func (r *Repository) GetFoodItem(ctx context.Context, id string) (*FoodItem, error) {
	var f FoodItem
	err := r.db.QueryRow(ctx,
		`SELECT id, organization_id, name, COALESCE(name_ne, ''), COALESCE(category, ''),
		        COALESCE(calories_per_100g, 0), COALESCE(protein_per_100g, 0),
		        COALESCE(carbs_per_100g, 0), COALESCE(fat_per_100g, 0),
		        COALESCE(fiber_per_100g, 0), COALESCE(serving_size_g, 0),
		        COALESCE(serving_label, ''), COALESCE(serving_label_ne, ''),
		        is_verified, COALESCE(barcode, ''), created_by, created_at, updated_at
		 FROM food_items
		 WHERE id = $1`,
		id,
	).Scan(&f.ID, &f.OrgID, &f.Name, &f.NameNe, &f.Category,
		&f.CaloriesPer100g, &f.ProteinPer100g, &f.CarbsPer100g, &f.FatPer100g,
		&f.FiberPer100g, &f.ServingSizeG, &f.ServingLabel, &f.ServingLabelNe,
		&f.IsVerified, &f.Barcode, &f.CreatedBy, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("getting food item: %w", err)
	}
	return &f, nil
}

// --- Food Logs ---

// CreateFoodLog inserts a new food log entry.
func (r *Repository) CreateFoodLog(ctx context.Context, userID, orgID, foodItemID, mealType string, quantityGrams, calories, protein, carbs, fat float64, loggedDate, notes string) (*FoodLog, error) {
	var l FoodLog
	err := r.db.QueryRow(ctx,
		`INSERT INTO food_logs (user_id, organization_id, food_item_id, meal_type, quantity_grams, calories, protein, carbs, fat, logged_date, notes)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		 RETURNING id, user_id, organization_id, food_item_id, meal_type,
		           quantity_grams, COALESCE(calories, 0), COALESCE(protein, 0),
		           COALESCE(carbs, 0), COALESCE(fat, 0),
		           logged_date, logged_at, COALESCE(notes, '')`,
		userID, nilIfEmpty(orgID), foodItemID, mealType, quantityGrams, calories, protein, carbs, fat, loggedDate, nilIfEmpty(notes),
	).Scan(&l.ID, &l.UserID, &l.OrgID, &l.FoodItemID, &l.MealType,
		&l.QuantityGrams, &l.Calories, &l.Protein, &l.Carbs, &l.Fat,
		&l.LoggedDate, &l.LoggedAt, &l.Notes)
	if err != nil {
		return nil, fmt.Errorf("creating food log: %w", err)
	}
	return &l, nil
}

// ListFoodLogs retrieves food logs for a user on a given date, optionally filtered by meal type.
func (r *Repository) ListFoodLogs(ctx context.Context, userID, orgID, date, mealType string) ([]FoodLog, error) {
	sql := `SELECT fl.id, fl.user_id, fl.organization_id, fl.food_item_id, fl.meal_type,
	               fl.quantity_grams, COALESCE(fl.calories, 0), COALESCE(fl.protein, 0),
	               COALESCE(fl.carbs, 0), COALESCE(fl.fat, 0),
	               fl.logged_date, fl.logged_at, COALESCE(fl.notes, ''),
	               fi.id, fi.organization_id, fi.name, COALESCE(fi.name_ne, ''), COALESCE(fi.category, ''),
	               COALESCE(fi.calories_per_100g, 0), COALESCE(fi.protein_per_100g, 0),
	               COALESCE(fi.carbs_per_100g, 0), COALESCE(fi.fat_per_100g, 0),
	               COALESCE(fi.fiber_per_100g, 0), COALESCE(fi.serving_size_g, 0),
	               COALESCE(fi.serving_label, ''), COALESCE(fi.serving_label_ne, ''),
	               fi.is_verified, COALESCE(fi.barcode, ''), fi.created_by, fi.created_at, fi.updated_at
	        FROM food_logs fl
	        JOIN food_items fi ON fi.id = fl.food_item_id
	        WHERE fl.user_id = $1 AND COALESCE(fl.organization_id, '00000000-0000-0000-0000-000000000000') = COALESCE($2, '00000000-0000-0000-0000-000000000000') AND fl.logged_date = $3`
	args := []interface{}{userID, nilIfEmpty(orgID), date}
	argIdx := 4

	if mealType != "" {
		sql += fmt.Sprintf(" AND fl.meal_type = $%d", argIdx)
		args = append(args, mealType)
	}

	sql += " ORDER BY fl.logged_at ASC"

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("listing food logs: %w", err)
	}
	defer rows.Close()

	var logs []FoodLog
	for rows.Next() {
		var l FoodLog
		var fi FoodItem
		if err := rows.Scan(&l.ID, &l.UserID, &l.OrgID, &l.FoodItemID, &l.MealType,
			&l.QuantityGrams, &l.Calories, &l.Protein, &l.Carbs, &l.Fat,
			&l.LoggedDate, &l.LoggedAt, &l.Notes,
			&fi.ID, &fi.OrgID, &fi.Name, &fi.NameNe, &fi.Category,
			&fi.CaloriesPer100g, &fi.ProteinPer100g, &fi.CarbsPer100g, &fi.FatPer100g,
			&fi.FiberPer100g, &fi.ServingSizeG, &fi.ServingLabel, &fi.ServingLabelNe,
			&fi.IsVerified, &fi.Barcode, &fi.CreatedBy, &fi.CreatedAt, &fi.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning food log: %w", err)
		}
		l.FoodItem = &fi
		logs = append(logs, l)
	}
	return logs, rows.Err()
}

// DeleteFoodLog deletes a food log entry owned by the user.
func (r *Repository) DeleteFoodLog(ctx context.Context, userID, logID string) error {
	tag, err := r.db.Exec(ctx,
		`DELETE FROM food_logs WHERE id = $1 AND user_id = $2`,
		logID, userID,
	)
	if err != nil {
		return fmt.Errorf("deleting food log: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("food log not found")
	}
	return nil
}

// --- Daily / Weekly Summary ---

// GetDailySummaryLogs retrieves all food logs for a user on a specific date, joined with food items.
func (r *Repository) GetDailySummaryLogs(ctx context.Context, userID, orgID, date string) ([]FoodLog, error) {
	rows, err := r.db.Query(ctx,
		`SELECT fl.id, fl.user_id, fl.organization_id, fl.food_item_id, fl.meal_type,
		        fl.quantity_grams, COALESCE(fl.calories, 0), COALESCE(fl.protein, 0),
		        COALESCE(fl.carbs, 0), COALESCE(fl.fat, 0),
		        fl.logged_date, fl.logged_at, COALESCE(fl.notes, ''),
		        fi.id, fi.organization_id, fi.name, COALESCE(fi.name_ne, ''), COALESCE(fi.category, ''),
		        COALESCE(fi.calories_per_100g, 0), COALESCE(fi.protein_per_100g, 0),
		        COALESCE(fi.carbs_per_100g, 0), COALESCE(fi.fat_per_100g, 0),
		        COALESCE(fi.fiber_per_100g, 0), COALESCE(fi.serving_size_g, 0),
		        COALESCE(fi.serving_label, ''), COALESCE(fi.serving_label_ne, ''),
		        fi.is_verified, COALESCE(fi.barcode, ''), fi.created_by, fi.created_at, fi.updated_at
		 FROM food_logs fl
		 JOIN food_items fi ON fi.id = fl.food_item_id
		 WHERE fl.user_id = $1 AND COALESCE(fl.organization_id, '00000000-0000-0000-0000-000000000000') = COALESCE($2, '00000000-0000-0000-0000-000000000000') AND fl.logged_date = $3
		 ORDER BY fl.meal_type, fl.logged_at ASC`,
		userID, nilIfEmpty(orgID), date,
	)
	if err != nil {
		return nil, fmt.Errorf("getting daily summary logs: %w", err)
	}
	defer rows.Close()

	var logs []FoodLog
	for rows.Next() {
		var l FoodLog
		var fi FoodItem
		if err := rows.Scan(&l.ID, &l.UserID, &l.OrgID, &l.FoodItemID, &l.MealType,
			&l.QuantityGrams, &l.Calories, &l.Protein, &l.Carbs, &l.Fat,
			&l.LoggedDate, &l.LoggedAt, &l.Notes,
			&fi.ID, &fi.OrgID, &fi.Name, &fi.NameNe, &fi.Category,
			&fi.CaloriesPer100g, &fi.ProteinPer100g, &fi.CarbsPer100g, &fi.FatPer100g,
			&fi.FiberPer100g, &fi.ServingSizeG, &fi.ServingLabel, &fi.ServingLabelNe,
			&fi.IsVerified, &fi.Barcode, &fi.CreatedBy, &fi.CreatedAt, &fi.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning daily summary log: %w", err)
		}
		l.FoodItem = &fi
		logs = append(logs, l)
	}
	return logs, rows.Err()
}

// GetWeeklySummary retrieves daily aggregated nutrition totals for 7 days starting from a date.
func (r *Repository) GetWeeklySummary(ctx context.Context, userID, orgID, fromDate string) ([]DailySummary, error) {
	rows, err := r.db.Query(ctx,
		`SELECT logged_date,
		        COALESCE(SUM(calories), 0),
		        COALESCE(SUM(protein), 0),
		        COALESCE(SUM(carbs), 0),
		        COALESCE(SUM(fat), 0)
		 FROM food_logs
		 WHERE user_id = $1 AND COALESCE(organization_id, '00000000-0000-0000-0000-000000000000') = COALESCE($2, '00000000-0000-0000-0000-000000000000')
		   AND logged_date >= $3::date AND logged_date < ($3::date + INTERVAL '7 days')
		 GROUP BY logged_date
		 ORDER BY logged_date ASC`,
		userID, nilIfEmpty(orgID), fromDate,
	)
	if err != nil {
		return nil, fmt.Errorf("getting weekly summary: %w", err)
	}
	defer rows.Close()

	var summaries []DailySummary
	for rows.Next() {
		var s DailySummary
		if err := rows.Scan(&s.Date, &s.TotalCalories, &s.TotalProtein, &s.TotalCarbs, &s.TotalFat); err != nil {
			return nil, fmt.Errorf("scanning weekly summary: %w", err)
		}
		summaries = append(summaries, s)
	}
	return summaries, rows.Err()
}

// --- Nutrition Goals ---

// UpsertNutritionGoal creates or updates the nutrition goal for a user+org pair.
func (r *Repository) UpsertNutritionGoal(ctx context.Context, userID, orgID string, calorieGoal, proteinGoalG, carbsGoalG, fatGoalG, weightKg, heightCm float64, age int, gender, activityLevel, goalType string) (*NutritionGoal, error) {
	var g NutritionGoal
	err := r.db.QueryRow(ctx,
		`INSERT INTO nutrition_goals (user_id, organization_id, calorie_goal, protein_goal_g, carbs_goal_g, fat_goal_g, weight_kg, height_cm, age, gender, activity_level, goal_type)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		 ON CONFLICT ON CONSTRAINT nutrition_goals_user_org_unique
		 DO UPDATE SET calorie_goal = $3, protein_goal_g = $4, carbs_goal_g = $5, fat_goal_g = $6,
		               weight_kg = $7, height_cm = $8, age = $9, gender = $10,
		               activity_level = $11, goal_type = $12
		 RETURNING id, user_id, organization_id, COALESCE(calorie_goal, 0), COALESCE(protein_goal_g, 0),
		           COALESCE(carbs_goal_g, 0), COALESCE(fat_goal_g, 0),
		           COALESCE(weight_kg, 0), COALESCE(height_cm, 0), COALESCE(age, 0),
		           COALESCE(gender, ''), COALESCE(activity_level, ''), COALESCE(goal_type, ''),
		           created_at, updated_at`,
		userID, nilIfEmpty(orgID), calorieGoal, proteinGoalG, carbsGoalG, fatGoalG, weightKg, heightCm, age, gender, activityLevel, goalType,
	).Scan(&g.ID, &g.UserID, &g.OrgID, &g.CalorieGoal, &g.ProteinGoalG,
		&g.CarbsGoalG, &g.FatGoalG, &g.WeightKg, &g.HeightCm, &g.Age,
		&g.Gender, &g.ActivityLevel, &g.GoalType, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("upserting nutrition goal: %w", err)
	}
	return &g, nil
}

// GetNutritionGoal retrieves the nutrition goal for a user+org pair.
func (r *Repository) GetNutritionGoal(ctx context.Context, userID, orgID string) (*NutritionGoal, error) {
	var g NutritionGoal
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, organization_id, COALESCE(calorie_goal, 0), COALESCE(protein_goal_g, 0),
		        COALESCE(carbs_goal_g, 0), COALESCE(fat_goal_g, 0),
		        COALESCE(weight_kg, 0), COALESCE(height_cm, 0), COALESCE(age, 0),
		        COALESCE(gender, ''), COALESCE(activity_level, ''), COALESCE(goal_type, ''),
		        created_at, updated_at
		 FROM nutrition_goals
		 WHERE user_id = $1 AND COALESCE(organization_id, '00000000-0000-0000-0000-000000000000') = COALESCE($2, '00000000-0000-0000-0000-000000000000')`,
		userID, nilIfEmpty(orgID),
	).Scan(&g.ID, &g.UserID, &g.OrgID, &g.CalorieGoal, &g.ProteinGoalG,
		&g.CarbsGoalG, &g.FatGoalG, &g.WeightKg, &g.HeightCm, &g.Age,
		&g.Gender, &g.ActivityLevel, &g.GoalType, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("getting nutrition goal: %w", err)
	}
	return &g, nil
}

// --- Meal Templates ---

// CreateMealTemplate inserts a new meal template.
func (r *Repository) CreateMealTemplate(ctx context.Context, userID, orgID, name, nameNe, mealType string, items json.RawMessage) (*MealTemplate, error) {
	var t MealTemplate
	err := r.db.QueryRow(ctx,
		`INSERT INTO meal_templates (user_id, organization_id, name, name_ne, meal_type, items)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, user_id, organization_id, name, COALESCE(name_ne, ''), COALESCE(meal_type, ''), items, created_at`,
		userID, nilIfEmpty(orgID), name, nilIfEmpty(nameNe), nilIfEmpty(mealType), items,
	).Scan(&t.ID, &t.UserID, &t.OrgID, &t.Name, &t.NameNe, &t.MealType, &t.Items, &t.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("creating meal template: %w", err)
	}
	return &t, nil
}

// ListMealTemplates retrieves meal templates for a user in an org.
func (r *Repository) ListMealTemplates(ctx context.Context, userID, orgID string) ([]MealTemplate, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, organization_id, name, COALESCE(name_ne, ''), COALESCE(meal_type, ''), items, created_at
		 FROM meal_templates
		 WHERE user_id = $1 AND COALESCE(organization_id, '00000000-0000-0000-0000-000000000000') = COALESCE($2, '00000000-0000-0000-0000-000000000000')
		 ORDER BY created_at DESC`,
		userID, nilIfEmpty(orgID),
	)
	if err != nil {
		return nil, fmt.Errorf("listing meal templates: %w", err)
	}
	defer rows.Close()

	var templates []MealTemplate
	for rows.Next() {
		var t MealTemplate
		if err := rows.Scan(&t.ID, &t.UserID, &t.OrgID, &t.Name, &t.NameNe, &t.MealType, &t.Items, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("scanning meal template: %w", err)
		}
		templates = append(templates, t)
	}
	return templates, rows.Err()
}

// GetMealTemplate retrieves a single meal template by ID.
func (r *Repository) GetMealTemplate(ctx context.Context, id string) (*MealTemplate, error) {
	var t MealTemplate
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, organization_id, name, COALESCE(name_ne, ''), COALESCE(meal_type, ''), items, created_at
		 FROM meal_templates
		 WHERE id = $1`,
		id,
	).Scan(&t.ID, &t.UserID, &t.OrgID, &t.Name, &t.NameNe, &t.MealType, &t.Items, &t.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("getting meal template: %w", err)
	}
	return &t, nil
}

// DeleteMealTemplate deletes a meal template owned by the user.
func (r *Repository) DeleteMealTemplate(ctx context.Context, userID, templateID string) error {
	tag, err := r.db.Exec(ctx,
		`DELETE FROM meal_templates WHERE id = $1 AND user_id = $2`,
		templateID, userID,
	)
	if err != nil {
		return fmt.Errorf("deleting meal template: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("meal template not found")
	}
	return nil
}

// --- Custom Foods ---

// CreateCustomFood inserts a user-created food item.
func (r *Repository) CreateCustomFood(ctx context.Context, userID string, input *CreateCustomFoodInput) (*FoodItem, error) {
	var f FoodItem
	err := r.db.QueryRow(ctx,
		`INSERT INTO food_items (name, name_ne, category, calories_per_100g, protein_per_100g,
		        carbs_per_100g, fat_per_100g, fiber_per_100g, serving_size_g, serving_label,
		        barcode, created_by, is_verified)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, false)
		 RETURNING id, organization_id, name, COALESCE(name_ne, ''), COALESCE(category, ''),
		           COALESCE(calories_per_100g, 0), COALESCE(protein_per_100g, 0),
		           COALESCE(carbs_per_100g, 0), COALESCE(fat_per_100g, 0),
		           COALESCE(fiber_per_100g, 0), COALESCE(serving_size_g, 0),
		           COALESCE(serving_label, ''), COALESCE(serving_label_ne, ''),
		           is_verified, COALESCE(barcode, ''), created_by, created_at, updated_at`,
		input.Name, nilIfEmpty(input.NameNe), nilIfEmpty(input.Category),
		input.CaloriesPer100g, input.ProteinPer100g, input.CarbsPer100g,
		input.FatPer100g, input.FiberPer100g, input.ServingSizeG,
		nilIfEmpty(input.ServingLabel), nilIfEmpty(input.Barcode), userID,
	).Scan(&f.ID, &f.OrgID, &f.Name, &f.NameNe, &f.Category,
		&f.CaloriesPer100g, &f.ProteinPer100g, &f.CarbsPer100g, &f.FatPer100g,
		&f.FiberPer100g, &f.ServingSizeG, &f.ServingLabel, &f.ServingLabelNe,
		&f.IsVerified, &f.Barcode, &f.CreatedBy, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("creating custom food: %w", err)
	}
	return &f, nil
}

// GetFoodByBarcode retrieves a food item by barcode.
func (r *Repository) GetFoodByBarcode(ctx context.Context, barcode string) (*FoodItem, error) {
	var f FoodItem
	err := r.db.QueryRow(ctx,
		`SELECT id, organization_id, name, COALESCE(name_ne, ''), COALESCE(category, ''),
		        COALESCE(calories_per_100g, 0), COALESCE(protein_per_100g, 0),
		        COALESCE(carbs_per_100g, 0), COALESCE(fat_per_100g, 0),
		        COALESCE(fiber_per_100g, 0), COALESCE(serving_size_g, 0),
		        COALESCE(serving_label, ''), COALESCE(serving_label_ne, ''),
		        is_verified, COALESCE(barcode, ''), created_by, created_at, updated_at
		 FROM food_items
		 WHERE barcode = $1`,
		barcode,
	).Scan(&f.ID, &f.OrgID, &f.Name, &f.NameNe, &f.Category,
		&f.CaloriesPer100g, &f.ProteinPer100g, &f.CarbsPer100g, &f.FatPer100g,
		&f.FiberPer100g, &f.ServingSizeG, &f.ServingLabel, &f.ServingLabelNe,
		&f.IsVerified, &f.Barcode, &f.CreatedBy, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("getting food by barcode: %w", err)
	}
	return &f, nil
}

// --- Nutrition Streaks ---

// UpdateNutritionStreak upserts the nutrition streak for a user based on the log date.
func (r *Repository) UpdateNutritionStreak(ctx context.Context, userID, logDate string) (*NutritionStreak, error) {
	var s NutritionStreak
	err := r.db.QueryRow(ctx,
		`INSERT INTO nutrition_streaks (user_id, current_streak, longest_streak, last_log_date, updated_at)
		 VALUES ($1, 1, 1, $2::date, NOW())
		 ON CONFLICT (user_id) DO UPDATE SET
		   current_streak = CASE
		     WHEN nutrition_streaks.last_log_date = $2::date THEN nutrition_streaks.current_streak
		     WHEN nutrition_streaks.last_log_date = $2::date - INTERVAL '1 day' THEN nutrition_streaks.current_streak + 1
		     ELSE 1
		   END,
		   longest_streak = GREATEST(
		     nutrition_streaks.longest_streak,
		     CASE
		       WHEN nutrition_streaks.last_log_date = $2::date THEN nutrition_streaks.current_streak
		       WHEN nutrition_streaks.last_log_date = $2::date - INTERVAL '1 day' THEN nutrition_streaks.current_streak + 1
		       ELSE 1
		     END
		   ),
		   last_log_date = $2::date,
		   updated_at = NOW()
		 RETURNING user_id, current_streak, longest_streak, last_log_date::text, updated_at`,
		userID, logDate,
	).Scan(&s.UserID, &s.CurrentStreak, &s.LongestStreak, &s.LastLogDate, &s.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("updating nutrition streak: %w", err)
	}
	return &s, nil
}

// GetNutritionStreak retrieves the nutrition streak for a user.
func (r *Repository) GetNutritionStreak(ctx context.Context, userID string) (*NutritionStreak, error) {
	var s NutritionStreak
	err := r.db.QueryRow(ctx,
		`SELECT user_id, current_streak, longest_streak, last_log_date::text, updated_at
		 FROM nutrition_streaks
		 WHERE user_id = $1`,
		userID,
	).Scan(&s.UserID, &s.CurrentStreak, &s.LongestStreak, &s.LastLogDate, &s.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("getting nutrition streak: %w", err)
	}
	return &s, nil
}

// --- Water Tracking ---

// GetWaterTotal returns the total water intake in ml for a user on a given date.
func (r *Repository) GetWaterTotal(ctx context.Context, userID, date string) (int, error) {
	var total int
	err := r.db.QueryRow(ctx,
		`SELECT COALESCE(SUM(amount_ml), 0) FROM water_logs WHERE user_id = $1 AND logged_date = $2`,
		userID, date,
	).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("getting water total: %w", err)
	}
	return total, nil
}

// GetUserWaterGoal returns the daily water goal in ml for a user.
func (r *Repository) GetUserWaterGoal(ctx context.Context, userID string) (int, error) {
	var goal int
	err := r.db.QueryRow(ctx,
		`SELECT COALESCE(daily_water_goal_ml, 2500) FROM users WHERE id = $1`,
		userID,
	).Scan(&goal)
	if err != nil {
		return 2500, fmt.Errorf("getting user water goal: %w", err)
	}
	return goal, nil
}

func nilIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
