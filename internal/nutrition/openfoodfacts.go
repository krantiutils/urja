package nutrition

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

var barcodePattern = regexp.MustCompile(`^\d{4,13}$`)

// offProduct holds parsed nutrition data from Open Food Facts.
type offProduct struct {
	Name            string
	CaloriesPer100g float64
	ProteinPer100g  float64
	CarbsPer100g    float64
	FatPer100g      float64
	FiberPer100g    float64
}

// offAPIResponse represents the Open Food Facts API response structure.
type offAPIResponse struct {
	Status  int `json:"status"`
	Product struct {
		ProductName string `json:"product_name"`
		Nutriments  struct {
			EnergyKcal100g   float64 `json:"energy-kcal_100g"`
			Proteins100g     float64 `json:"proteins_100g"`
			Carbohydrates100g float64 `json:"carbohydrates_100g"`
			Fat100g          float64 `json:"fat_100g"`
			Fiber100g        float64 `json:"fiber_100g"`
		} `json:"nutriments"`
	} `json:"product"`
}

var offHTTPClient = &http.Client{Timeout: 10 * time.Second}

// fetchOpenFoodFacts fetches food data from the Open Food Facts API by barcode.
func fetchOpenFoodFacts(barcode string) (*offProduct, error) {
	if !barcodePattern.MatchString(barcode) {
		return nil, fmt.Errorf("invalid barcode format")
	}
	url := fmt.Sprintf("https://world.openfoodfacts.org/api/v0/product/%s.json", barcode)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating OFF request: %w", err)
	}
	req.Header.Set("User-Agent", "UrjaFitnessApp/1.0")

	resp, err := offHTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching from OFF: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OFF returned status %d", resp.StatusCode)
	}

	var apiResp offAPIResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("decoding OFF response: %w", err)
	}

	if apiResp.Status != 1 || apiResp.Product.ProductName == "" {
		return nil, fmt.Errorf("product not found on OFF")
	}

	return &offProduct{
		Name:            apiResp.Product.ProductName,
		CaloriesPer100g: apiResp.Product.Nutriments.EnergyKcal100g,
		ProteinPer100g:  apiResp.Product.Nutriments.Proteins100g,
		CarbsPer100g:    apiResp.Product.Nutriments.Carbohydrates100g,
		FatPer100g:      apiResp.Product.Nutriments.Fat100g,
		FiberPer100g:    apiResp.Product.Nutriments.Fiber100g,
	}, nil
}
