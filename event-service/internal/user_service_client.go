package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/anrisys/quicket/event-service/pkg/config"
	"github.com/anrisys/quicket/event-service/pkg/errs"
)

type UserReader interface {
	GetUserID(ctx context.Context, publicID string) (*uint, error)
	FindUserByPublicID(ctx context.Context, publicID string) (*UserDTO, error)
}


type UserServiceClient struct {
	cfg    *config.Config
	httpClient *http.Client
}

func NewUserServiceClient(cfg *config.Config) *UserServiceClient  {
	return &UserServiceClient{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *UserServiceClient) GetUserID(ctx context.Context, publicID string) (*uint, error) {
	baseURL := c.cfg.Clients.UserServiceURL
	url := fmt.Sprintf("%s/api/v1/users/%s/primary-id", baseURL, publicID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("user client#GetUserID: failed to create request: %w", err)
	}

	if authToken, ok := ctx.Value("Authorization").(string); ok {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user client#GetUserID: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errs.NewErrNotFound("user")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user client#GetUserID: unexpected status: %d", resp.StatusCode)
	}

	var res struct {
		Code string
		Message string
		PrimaryID uint
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("user client#GetUserID: failed to decode response: %w", err)
	}

	return &res.PrimaryID, nil
}

func (c *UserServiceClient) FindUserByPublicID(ctx context.Context, publicID string) (*UserDTO, error){
	baseURL := c.cfg.Clients.UserServiceURL
	url := fmt.Sprintf("%s/api/v1/users/public/%s", baseURL, publicID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("user_client#FindUserByPulbicID: failed to create request: %w", err)
	}

	if authToken, ok := ctx.Value("Authorization").(string); ok {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user_client#FindUserByPulbicID: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errs.NewErrNotFound("user")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user_client#FindUserByPulbicID: unexpected status: %d", resp.StatusCode)
	}

	var res *UserDTO
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("user_client#FindUserByPulbicID: failed to decode response: %w", err)
	}

	return res, nil
}

