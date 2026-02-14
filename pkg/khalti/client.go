package khalti

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// Client communicates with the Khalti e-payment API.
type Client struct {
	secretKey  string
	baseURL    string
	httpClient *http.Client
	logger     *slog.Logger
}

// NewClient creates a new Khalti API client.
func NewClient(secretKey, baseURL string, logger *slog.Logger) *Client {
	return &Client{
		secretKey: secretKey,
		baseURL:   baseURL,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		logger: logger,
	}
}

// LookupResponse is the Khalti payment verification response.
type LookupResponse struct {
	Pidx          string `json:"pidx"`
	TotalAmount   int    `json:"total_amount"` // in paisa (NPR * 100)
	Status        string `json:"status"`       // Completed, Pending, Initiated, Refunded, Expired
	TransactionID string `json:"transaction_id"`
	Fee           int    `json:"fee"`
	Refunded      bool   `json:"refunded"`
}

// LookupPayment verifies a Khalti payment by its pidx.
// Returns the payment details or an error.
func (c *Client) LookupPayment(ctx context.Context, pidx string) (*LookupResponse, error) {
	if c.secretKey == "" {
		return nil, fmt.Errorf("khalti secret key not configured")
	}

	body, err := json.Marshal(map[string]string{"pidx": pidx})
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	url := c.baseURL + "/api/v2/epayment/lookup/"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "key "+c.secretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("khalti lookup request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading khalti response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("khalti lookup failed",
			"status", resp.StatusCode,
			"body", string(respBody),
			"pidx", pidx,
		)
		return nil, fmt.Errorf("khalti lookup failed with status %d", resp.StatusCode)
	}

	var result LookupResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("decoding khalti response: %w", err)
	}

	return &result, nil
}
