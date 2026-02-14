package health

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/urja-gym/urja/pkg/middleware"
)

const maxUploadSize = 10 << 20 // 10 MB

// Handler handles HTTP requests for health endpoints.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new health handler.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// GetMetrics handles GET /api/v1/members/me/health
func (h *Handler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	q := r.URL.Query()
	metricType := q.Get("type")
	limitStr := q.Get("limit")

	var from, to *time.Time
	if fromStr := q.Get("from"); fromStr != "" {
		t, err := time.Parse("2006-01-02", fromStr)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid 'from' date format, use YYYY-MM-DD"})
			return
		}
		from = &t
	}
	if toStr := q.Get("to"); toStr != "" {
		t, err := time.Parse("2006-01-02", toStr)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid 'to' date format, use YYYY-MM-DD"})
			return
		}
		// Include the full "to" day
		endOfDay := t.Add(24*time.Hour - time.Nanosecond)
		to = &endOfDay
	}

	limit, _ := strconv.Atoi(limitStr)

	metrics, err := h.service.GetMetrics(r.Context(), userID, metricType, from, to, limit)
	if err != nil {
		if isValidationError(err) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		h.logger.Error("failed to get health metrics", "error", err, "user_id", userID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data": metrics,
	})
}

// LogBMI handles POST /api/v1/members/me/health/bmi
func (h *Handler) LogBMI(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var input BMIInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	metric, err := h.service.LogBMI(r.Context(), userID, &input)
	if err != nil {
		if isValidationError(err) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		h.logger.Error("failed to log BMI", "error", err, "user_id", userID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusCreated, metric)
}

// LogMeasurements handles POST /api/v1/members/me/health/measurements
func (h *Handler) LogMeasurements(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var input MeasurementsInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	metric, err := h.service.LogMeasurements(r.Context(), userID, &input)
	if err != nil {
		if isValidationError(err) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		h.logger.Error("failed to log measurements", "error", err, "user_id", userID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusCreated, metric)
}

// UploadPhoto handles POST /api/v1/members/me/health/photos
func (h *Handler) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "file too large or invalid multipart form"})
		return
	}

	file, header, err := r.FormFile("photo")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "photo file is required"})
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	input := &PhotoInput{
		Category: r.FormValue("category"),
		Notes:    r.FormValue("notes"),
		FileExt:  ext,
		FileData: file,
	}

	photo, err := h.service.UploadPhoto(r.Context(), userID, input)
	if err != nil {
		if isValidationError(err) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		h.logger.Error("failed to upload photo", "error", err, "user_id", userID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusCreated, photo)
}

// ListPhotos handles GET /api/v1/members/me/health/photos
func (h *Handler) ListPhotos(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	q := r.URL.Query()
	category := q.Get("category")
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))

	photos, total, err := h.service.ListPhotos(r.Context(), userID, category, limit, offset)
	if err != nil {
		if isValidationError(err) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		h.logger.Error("failed to list photos", "error", err, "user_id", userID)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	if photos == nil {
		photos = []ProgressPhoto{}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":  photos,
		"total": total,
	})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}

func isValidationError(err error) bool {
	unwrapper, ok := err.(interface{ Unwrap() error })
	if !ok {
		return true
	}
	return unwrapper.Unwrap() == nil
}
