package nutrition

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/urja-gym/urja/pkg/middleware"
)

// Handler handles HTTP requests for nutrition endpoints.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new nutrition handler.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// --- Food Items ---

// SearchFoods handles GET /api/v1/members/me/foods
func (h *Handler) SearchFoods(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	orgID := r.URL.Query().Get("organization_id")
	query := r.URL.Query().Get("q")
	category := r.URL.Query().Get("category")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	items, err := h.service.SearchFoods(r.Context(), userID, orgID, query, category, limit)
	if err != nil {
		h.logger.Error("failed to search foods", "error", err, "org_id", orgID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": items})
}

// GetFoodItem handles GET /api/v1/members/me/foods/{id}
func (h *Handler) GetFoodItem(w http.ResponseWriter, r *http.Request) {
	foodID := chi.URLParam(r, "id")
	if foodID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing food item ID"})
		return
	}

	item, err := h.service.GetFoodItem(r.Context(), foodID)
	if err != nil {
		h.logger.Error("failed to get food item", "error", err, "food_id", foodID)
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "food item not found"})
		return
	}

	writeJSON(w, http.StatusOK, item)
}

// CreateCustomFood handles POST /api/v1/members/me/foods
func (h *Handler) CreateCustomFood(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var input CreateCustomFoodInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	f, err := h.service.CreateCustomFood(r.Context(), userID, &input)
	if err != nil {
		h.logger.Error("failed to create custom food", "error", err, "user_id", userID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, f)
}

// LookupBarcode handles GET /api/v1/members/me/foods/barcode/{barcode}
func (h *Handler) LookupBarcode(w http.ResponseWriter, r *http.Request) {
	barcode := chi.URLParam(r, "barcode")
	if barcode == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing barcode"})
		return
	}

	item, err := h.service.LookupBarcode(r.Context(), barcode)
	if err != nil {
		h.logger.Error("failed to lookup barcode", "error", err, "barcode", barcode)
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, item)
}

// GetNutritionStreak handles GET /api/v1/members/me/nutrition/streak
func (h *Handler) GetNutritionStreak(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	streak, err := h.service.GetNutritionStreak(r.Context(), userID)
	if err != nil {
		h.logger.Error("failed to get nutrition streak", "error", err, "user_id", userID)
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "nutrition streak not found"})
		return
	}

	writeJSON(w, http.StatusOK, streak)
}

// --- Food Logs ---

// CreateFoodLog handles POST /api/v1/members/me/food-logs
func (h *Handler) CreateFoodLog(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var input CreateFoodLogInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	l, err := h.service.CreateFoodLog(r.Context(), userID, &input)
	if err != nil {
		h.logger.Error("failed to create food log", "error", err, "user_id", userID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, l)
}

// ListFoodLogs handles GET /api/v1/members/me/food-logs
func (h *Handler) ListFoodLogs(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	orgID := r.URL.Query().Get("organization_id")
	date := r.URL.Query().Get("date")
	mealType := r.URL.Query().Get("meal_type")

	logs, err := h.service.ListFoodLogs(r.Context(), userID, orgID, date, mealType)
	if err != nil {
		h.logger.Error("failed to list food logs", "error", err, "user_id", userID, "org_id", orgID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": logs})
}

// DeleteFoodLog handles DELETE /api/v1/members/me/food-logs/{id}
func (h *Handler) DeleteFoodLog(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	logID := chi.URLParam(r, "id")
	if logID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing food log ID"})
		return
	}

	if err := h.service.DeleteFoodLog(r.Context(), userID, logID); err != nil {
		h.logger.Error("failed to delete food log", "error", err, "log_id", logID, "user_id", userID)
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "food log deleted"})
}

// --- Daily / Weekly Summary ---

// GetDailySummary handles GET /api/v1/members/me/nutrition/summary
func (h *Handler) GetDailySummary(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	orgID := r.URL.Query().Get("organization_id")
	date := r.URL.Query().Get("date")

	summary, err := h.service.GetDailySummary(r.Context(), userID, orgID, date)
	if err != nil {
		h.logger.Error("failed to get daily summary", "error", err, "user_id", userID, "org_id", orgID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, summary)
}

// GetWeeklySummary handles GET /api/v1/members/me/nutrition/weekly
func (h *Handler) GetWeeklySummary(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	orgID := r.URL.Query().Get("organization_id")
	fromDate := r.URL.Query().Get("from")

	summaries, err := h.service.GetWeeklySummary(r.Context(), userID, orgID, fromDate)
	if err != nil {
		h.logger.Error("failed to get weekly summary", "error", err, "user_id", userID, "org_id", orgID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": summaries})
}

// GetDailyDashboard handles GET /api/v1/members/me/nutrition/daily-dashboard
func (h *Handler) GetDailyDashboard(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	orgID := r.URL.Query().Get("organization_id")
	date := r.URL.Query().Get("date")

	dashboard, err := h.service.GetDailyDashboard(r.Context(), userID, orgID, date)
	if err != nil {
		h.logger.Error("failed to get daily dashboard", "error", err, "user_id", userID, "org_id", orgID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, dashboard)
}

// --- Nutrition Goals ---

// SetNutritionGoal handles POST /api/v1/members/me/nutrition/goal
func (h *Handler) SetNutritionGoal(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var input SetNutritionGoalInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	g, err := h.service.SetNutritionGoal(r.Context(), userID, &input)
	if err != nil {
		h.logger.Error("failed to set nutrition goal", "error", err, "user_id", userID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, g)
}

// GetNutritionGoal handles GET /api/v1/members/me/nutrition/goal
func (h *Handler) GetNutritionGoal(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	orgID := r.URL.Query().Get("organization_id")

	g, err := h.service.GetNutritionGoal(r.Context(), userID, orgID)
	if err != nil {
		h.logger.Error("failed to get nutrition goal", "error", err, "user_id", userID, "org_id", orgID)
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "nutrition goal not found"})
		return
	}

	writeJSON(w, http.StatusOK, g)
}

// --- Meal Templates ---

// CreateMealTemplate handles POST /api/v1/members/me/meal-templates
func (h *Handler) CreateMealTemplate(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var input CreateMealTemplateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	t, err := h.service.CreateMealTemplate(r.Context(), userID, &input)
	if err != nil {
		h.logger.Error("failed to create meal template", "error", err, "user_id", userID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, t)
}

// ListMealTemplates handles GET /api/v1/members/me/meal-templates
func (h *Handler) ListMealTemplates(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	orgID := r.URL.Query().Get("organization_id")

	templates, err := h.service.ListMealTemplates(r.Context(), userID, orgID)
	if err != nil {
		h.logger.Error("failed to list meal templates", "error", err, "user_id", userID, "org_id", orgID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"data": templates})
}

// DeleteMealTemplate handles DELETE /api/v1/members/me/meal-templates/{id}
func (h *Handler) DeleteMealTemplate(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	templateID := chi.URLParam(r, "id")
	if templateID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing meal template ID"})
		return
	}

	if err := h.service.DeleteMealTemplate(r.Context(), userID, templateID); err != nil {
		h.logger.Error("failed to delete meal template", "error", err, "template_id", templateID, "user_id", userID)
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "meal template deleted"})
}

// LogMealTemplate handles POST /api/v1/members/me/meal-templates/{id}/log
func (h *Handler) LogMealTemplate(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	templateID := chi.URLParam(r, "id")
	if templateID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing meal template ID"})
		return
	}

	var input LogMealTemplateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	logs, err := h.service.LogMealTemplate(r.Context(), userID, templateID, &input)
	if err != nil {
		h.logger.Error("failed to log meal template", "error", err, "template_id", templateID, "user_id", userID)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{"data": logs})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}
