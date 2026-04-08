package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"quicket-payment-service/internal/helper"
	"quicket-payment-service/internal/payment/dto"
	"time"
)

type HTTPPaymentGateway struct {
	client  *http.Client
	baseURL string
	apiKey  string
}

func NewHTTPPaymentGateway(baseURL, apiKey string, timeout time.Duration) *HTTPPaymentGateway {
	return &HTTPPaymentGateway{
		client:  helper.NewHTTPClient(timeout),
		baseURL: baseURL,
		apiKey:  apiKey,
	}
}

func (g *HTTPPaymentGateway) CreatePayment(ctx context.Context, req dto.CreatePaymentRequest) (*dto.CreatePaymentResponse, error) {

	url := g.baseURL + "/payment_requests"

	payload := map[string]any{
		"external_id": req.ExternalID,
		"amount":      req.Amount,
		"currency":    req.Currency,
		"expiry_date": req.Expiry.Format(time.RFC3339),
	}

	headers := map[string]string{
		"Authorization":   "Bearer " + g.apiKey,
		"Idempotency-Key": req.ExternalID,
	}

	resp, err := helper.RetryWithJitter(ctx, 3, 100*time.Millisecond, func() (*http.Response, error) {
		return helper.DoJSONRequest(ctx, g.client, http.MethodPost, url, headers, payload)
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("gateway error: %s", string(body))
	}

	var result struct {
		ID         string `json:"id"`
		InvoiceURL string `json:"invoice_url"`
		Status     string `json:"status"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &dto.CreatePaymentResponse{
		GatewayID:  result.ID,
		PaymentURL: result.InvoiceURL,
		Status:     result.Status,
	}, nil
}
