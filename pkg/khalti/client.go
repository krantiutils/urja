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

// Client communicates with the Khalti ePayment API.
type Client struct {
	secretKey  string
	baseURL    string
	websiteURL string
	returnURL  string
	httpClient *http.Client
	logger     *slog.Logger
}

// NewClient creates a Khalti API client.
func NewClient(secretKey, baseURL, websiteURL, returnURL string, logger *slog.Logger) *Client {
	return &Client{
		secretKey:  secretKey,
		baseURL:    baseURL,
		websiteURL: websiteURL,
		returnURL:  returnURL,
		httpClient: &http.Client{Timeout: 15 * time.Second},
		logger:     logger,
	}
}

// InitiateRequest is the payload for initiating a Khalti payment.
type InitiateRequest struct {
	ReturnURL         string        `json:"return_url"`
	WebsiteURL        string        `json:"website_url"`
	Amount            int64         `json:"amount"` // in paisa
	PurchaseOrderID   string        `json:"purchase_order_id"`
	PurchaseOrderName string        `json:"purchase_order_name"`
	CustomerInfo      *CustomerInfo `json:"customer_info,omitempty"`
}

// CustomerInfo is optional customer metadata for Khalti.
type CustomerInfo struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

// InitiateResponse is the response from Khalti payment initiation.
type InitiateResponse struct {
	Pidx       string `json:"pidx"`
	PaymentURL string `json:"payment_url"`
	ExpiresAt  string `json:"expires_at"`
	ExpiresIn  int    `json:"expires_in"`
}

// Initiate starts a Khalti ePayment session.
func (c *Client) Initiate(ctx context.Context, orderID, orderName string, amountPaisa int64, customer *CustomerInfo) (*InitiateResponse, error) {
	req := InitiateRequest{
		ReturnURL:         c.returnURL,
		WebsiteURL:        c.websiteURL,
		Amount:            amountPaisa,
		PurchaseOrderID:   orderID,
		PurchaseOrderName: orderName,
		CustomerInfo:      customer,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshaling initiate request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/epayment/initiate/", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating initiate request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Key "+c.secretKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("khalti initiate request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading initiate response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("khalti initiate failed", "status", resp.StatusCode, "body", string(respBody))
		return nil, fmt.Errorf("khalti initiate returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result InitiateResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("decoding initiate response: %w", err)
	}

	c.logger.Info("khalti payment initiated", "pidx", result.Pidx, "order_id", orderID)
	return &result, nil
}

// LookupResponse is the response from Khalti payment lookup.
type LookupResponse struct {
	Pidx          string `json:"pidx"`
	TotalAmount   int64  `json:"total_amount"`
	Status        string `json:"status"`
	TransactionID string `json:"transaction_id"`
	Fee           int64  `json:"fee"`
	Refunded      bool   `json:"refunded"`
}

// Lookup verifies a Khalti payment by pidx.
func (c *Client) Lookup(ctx context.Context, pidx string) (*LookupResponse, error) {
	body, err := json.Marshal(map[string]string{"pidx": pidx})
	if err != nil {
		return nil, fmt.Errorf("marshaling lookup request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/epayment/lookup/", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating lookup request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Key "+c.secretKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("khalti lookup request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading lookup response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("khalti lookup failed", "status", resp.StatusCode, "body", string(respBody))
		return nil, fmt.Errorf("khalti lookup returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result LookupResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("decoding lookup response: %w", err)
	}

	c.logger.Info("khalti lookup completed", "pidx", pidx, "status", result.Status)
	return &result, nil
}

// IsCompleted returns true if the lookup status indicates successful payment.
func (r *LookupResponse) IsCompleted() bool {
	return r.Status == "Completed"
}
