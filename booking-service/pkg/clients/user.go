package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"quicket/booking-service/pkg/config"
	"time"

	"github.com/rs/zerolog"
)

type UserServiceClient struct {
	cfg        *config.ClientServices
	httpClient *http.Client
	logger zerolog.Logger
}

func NewUserServiceClient(cfg *config.Config, logger zerolog.Logger) *UserServiceClient {
	return &UserServiceClient{
		cfg: cfg.Clients,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		logger: logger,
	}
}

func (c *UserServiceClient) GetUserID(ctx context.Context, publicID string) (*uint, error) {
	log := c.logger.With().
		Str("client_services", "user_client").
		Str("method", "get_user_id").
		Str("user_public_id", publicID).
		Logger()
	
	baseURL := c.cfg.UserServiceURL
	url := fmt.Sprintf("%s/api/v1/users/%s/primary-id", baseURL, publicID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to create request")
		return nil, fmt.Errorf("user client#GetUserID: failed to create request: %w", err)
	}

	if authToken, ok := ctx.Value("Authorization").(string); ok {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("http request failed")
		return nil, fmt.Errorf("%w: %v", ErrRequestClientFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrServiceClientUserNotFound
	}

	if resp.StatusCode != http.StatusOK {
		log.Error().Err(err).Int("response_status", resp.StatusCode).Msg("unexpected response status")
		return nil, fmt.Errorf("%w: %d", ErrClientUnexpectedResponseStatus, resp.StatusCode)
	}

	var res struct {
		Code string
		Message string
		PrimaryID uint
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Error().Err(err).Msg("failed to decode response")
		return nil, fmt.Errorf("%w: %v", ErrClientFailedToDecodeResponse, err)
	}

	return &res.PrimaryID, nil
}


