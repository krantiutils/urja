package nutrition

import (
	"github.com/go-chi/chi/v5"
)

// RegisterSelfRoutes mounts member self-service nutrition routes (under /members/me).
func (h *Handler) RegisterSelfRoutes(r chi.Router) {
	// Food search and custom foods
	r.Get("/foods", h.SearchFoods)
	r.Post("/foods", h.CreateCustomFood)
	r.Get("/foods/barcode/{barcode}", h.LookupBarcode)
	r.Get("/foods/{id}", h.GetFoodItem)

	// Food logs
	r.Post("/food-logs", h.CreateFoodLog)
	r.Get("/food-logs", h.ListFoodLogs)
	r.Delete("/food-logs/{id}", h.DeleteFoodLog)

	// Nutrition summaries and goals
	r.Get("/nutrition/summary", h.GetDailySummary)
	r.Get("/nutrition/weekly", h.GetWeeklySummary)
	r.Get("/nutrition/daily-dashboard", h.GetDailyDashboard)
	r.Post("/nutrition/goal", h.SetNutritionGoal)
	r.Get("/nutrition/goal", h.GetNutritionGoal)
	r.Get("/nutrition/streak", h.GetNutritionStreak)

	// Meal templates
	r.Post("/meal-templates", h.CreateMealTemplate)
	r.Get("/meal-templates", h.ListMealTemplates)
	r.Delete("/meal-templates/{id}", h.DeleteMealTemplate)
	r.Post("/meal-templates/{id}/log", h.LogMealTemplate)
}
