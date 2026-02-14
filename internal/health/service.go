package health

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// BMI input from the client.
type BMIInput struct {
	HeightCM float64 `json:"height_cm"`
	WeightKG float64 `json:"weight_kg"`
}

// BMIValue is the JSONB value stored for a BMI metric.
type BMIValue struct {
	HeightCM float64 `json:"height_cm"`
	WeightKG float64 `json:"weight_kg"`
	BMI      float64 `json:"bmi"`
	Category string  `json:"category"`
}

// MeasurementsInput from the client.
type MeasurementsInput struct {
	ChestCM  *float64 `json:"chest_cm,omitempty"`
	WaistCM  *float64 `json:"waist_cm,omitempty"`
	HipsCM   *float64 `json:"hips_cm,omitempty"`
	ArmsCM   *float64 `json:"arms_cm,omitempty"`
	ThighsCM *float64 `json:"thighs_cm,omitempty"`
}

// MeasurementsValue is the JSONB value stored for a measurements metric.
type MeasurementsValue struct {
	ChestCM  *float64 `json:"chest_cm,omitempty"`
	WaistCM  *float64 `json:"waist_cm,omitempty"`
	HipsCM   *float64 `json:"hips_cm,omitempty"`
	ArmsCM   *float64 `json:"arms_cm,omitempty"`
	ThighsCM *float64 `json:"thighs_cm,omitempty"`
	WHR      *float64 `json:"whr,omitempty"`
}

// PhotoInput holds metadata for a progress photo upload.
type PhotoInput struct {
	Category string
	Notes    string
	FileExt  string
	FileData io.Reader
}

// Service handles health business logic.
type Service struct {
	repo     *Repository
	uploadDir string
	logger   *slog.Logger
}

// NewService creates a new health service.
func NewService(repo *Repository, uploadDir string, logger *slog.Logger) *Service {
	return &Service{repo: repo, uploadDir: uploadDir, logger: logger}
}

// LogBMI records a BMI metric with auto-calculated BMI and category.
func (s *Service) LogBMI(ctx context.Context, memberID string, input *BMIInput) (*HealthMetric, error) {
	if err := validateBMIInput(input); err != nil {
		return nil, err
	}

	bmi := input.WeightKG / math.Pow(input.HeightCM/100.0, 2)
	bmi = math.Round(bmi*100) / 100

	val := BMIValue{
		HeightCM: input.HeightCM,
		WeightKG: input.WeightKG,
		BMI:      bmi,
		Category: bmiCategory(bmi),
	}

	valueJSON, err := json.Marshal(val)
	if err != nil {
		return nil, fmt.Errorf("marshaling BMI value: %w", err)
	}

	metric, err := s.repo.InsertMetric(ctx, memberID, "bmi", valueJSON)
	if err != nil {
		return nil, fmt.Errorf("logging BMI: %w", err)
	}
	return metric, nil
}

// LogMeasurements records body measurements with auto-calculated WHR.
func (s *Service) LogMeasurements(ctx context.Context, memberID string, input *MeasurementsInput) (*HealthMetric, error) {
	if err := validateMeasurementsInput(input); err != nil {
		return nil, err
	}

	val := MeasurementsValue{
		ChestCM:  input.ChestCM,
		WaistCM:  input.WaistCM,
		HipsCM:   input.HipsCM,
		ArmsCM:   input.ArmsCM,
		ThighsCM: input.ThighsCM,
	}

	// Auto-calculate WHR if both waist and hips are provided
	if input.WaistCM != nil && input.HipsCM != nil && *input.HipsCM > 0 {
		whr := *input.WaistCM / *input.HipsCM
		whr = math.Round(whr*100) / 100
		val.WHR = &whr
	}

	valueJSON, err := json.Marshal(val)
	if err != nil {
		return nil, fmt.Errorf("marshaling measurements value: %w", err)
	}

	metric, err := s.repo.InsertMetric(ctx, memberID, "measurements", valueJSON)
	if err != nil {
		return nil, fmt.Errorf("logging measurements: %w", err)
	}
	return metric, nil
}

// GetMetrics retrieves health metrics with optional filtering.
func (s *Service) GetMetrics(ctx context.Context, memberID, metricType string, from, to *time.Time, limit int) ([]HealthMetric, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}

	validTypes := map[string]bool{
		"bmi": true, "whr": true, "weight": true,
		"body_fat": true, "measurements": true,
	}
	if metricType != "" && !validTypes[metricType] {
		return nil, errors.New("invalid metric_type: must be bmi, whr, weight, body_fat, or measurements")
	}

	metrics, err := s.repo.ListMetrics(ctx, memberID, metricType, from, to, limit)
	if err != nil {
		return nil, fmt.Errorf("getting metrics: %w", err)
	}
	if metrics == nil {
		metrics = []HealthMetric{}
	}
	return metrics, nil
}

// UploadPhoto saves a progress photo to disk and records metadata.
func (s *Service) UploadPhoto(ctx context.Context, memberID string, input *PhotoInput) (*ProgressPhoto, error) {
	if err := validatePhotoInput(input); err != nil {
		return nil, err
	}

	// Create member-specific upload directory
	memberDir := filepath.Join(s.uploadDir, memberID)
	if err := os.MkdirAll(memberDir, 0755); err != nil {
		return nil, fmt.Errorf("creating upload directory: %w", err)
	}

	// Generate unique filename
	filename := uuid.New().String() + input.FileExt
	filePath := filepath.Join(memberDir, filename)

	f, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("creating photo file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, input.FileData); err != nil {
		os.Remove(filePath)
		return nil, fmt.Errorf("writing photo file: %w", err)
	}

	// Store relative path as the photo URL
	photoURL := filepath.Join(memberID, filename)

	photo, err := s.repo.InsertPhoto(ctx, memberID, photoURL, input.Category, input.Notes)
	if err != nil {
		os.Remove(filePath)
		return nil, fmt.Errorf("saving photo metadata: %w", err)
	}
	return photo, nil
}

// ListPhotos retrieves progress photos with optional category filter.
func (s *Service) ListPhotos(ctx context.Context, memberID, category string, limit, offset int) ([]ProgressPhoto, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	validCategories := map[string]bool{"front": true, "side": true, "back": true}
	if category != "" && !validCategories[category] {
		return nil, 0, errors.New("invalid category: must be front, side, or back")
	}

	return s.repo.ListPhotos(ctx, memberID, category, limit, offset)
}

func validateBMIInput(input *BMIInput) error {
	if input.HeightCM <= 0 || input.HeightCM > 300 {
		return errors.New("height_cm must be between 0 and 300")
	}
	if input.WeightKG <= 0 || input.WeightKG > 700 {
		return errors.New("weight_kg must be between 0 and 700")
	}
	return nil
}

func validateMeasurementsInput(input *MeasurementsInput) error {
	if input.ChestCM == nil && input.WaistCM == nil && input.HipsCM == nil &&
		input.ArmsCM == nil && input.ThighsCM == nil {
		return errors.New("at least one measurement is required")
	}

	check := func(name string, val *float64) error {
		if val != nil && (*val <= 0 || *val > 300) {
			return fmt.Errorf("%s must be between 0 and 300 cm", name)
		}
		return nil
	}

	if err := check("chest_cm", input.ChestCM); err != nil {
		return err
	}
	if err := check("waist_cm", input.WaistCM); err != nil {
		return err
	}
	if err := check("hips_cm", input.HipsCM); err != nil {
		return err
	}
	if err := check("arms_cm", input.ArmsCM); err != nil {
		return err
	}
	return check("thighs_cm", input.ThighsCM)
}

func validatePhotoInput(input *PhotoInput) error {
	if input.FileData == nil {
		return errors.New("photo file is required")
	}

	validCategories := map[string]bool{"front": true, "side": true, "back": true, "": true}
	if !validCategories[input.Category] {
		return errors.New("invalid category: must be front, side, or back")
	}

	validExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !validExts[strings.ToLower(input.FileExt)] {
		return errors.New("invalid file type: must be jpg, jpeg, png, or webp")
	}

	return nil
}

func bmiCategory(bmi float64) string {
	switch {
	case bmi < 18.5:
		return "underweight"
	case bmi < 25.0:
		return "normal"
	case bmi < 30.0:
		return "overweight"
	default:
		return "obese"
	}
}
