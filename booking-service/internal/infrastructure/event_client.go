package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"quicket/booking-service/internal/domain/booking"
	"time"
)

type HTTPEventClient struct {
	baseURL string
	client  *http.Client
}

func NewHTTPEventClient(baseURL string) *HTTPEventClient {
	return &HTTPEventClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 5 * time.Second},
	}
}
func (c *HTTPEventClient) ReserveSeats(
	ctx context.Context,
	eventPublicID string,
	seats uint64,
) (uint64, error) {

	url := fmt.Sprintf("%s/internal/events/%s/reserve", c.baseURL, eventPublicID)

	payload := map[string]uint64{
		"seats": seats,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}

	var seatPrice uint64

	err = Retry(ctx, 3, 100*time.Millisecond, func() error {

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			url,
			bytes.NewReader(body),
		)
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := c.client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// SUCCESS RESPONSE
		if resp.StatusCode == http.StatusOK {

			var successResp reserveSeatsSuccessResponse

			if err := json.NewDecoder(resp.Body).Decode(&successResp); err != nil {
				return booking.ErrEventServiceUnavailable
			}

			seatPrice = successResp.SeatPrice
			return nil
		}

		// BUSINESS ERROR → do NOT retry
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {

			var svcErr clientServiceError
			if err := json.NewDecoder(resp.Body).Decode(&svcErr); err != nil {
				return booking.ErrEventServiceUnavailable
			}

			return mapEventServiceError(svcErr.Code)
		}

		// 5xx → retry
		return booking.ErrEventServiceUnavailable
	})

	if err != nil {
		return 0, err
	}

	return seatPrice, nil
}

func (c *HTTPEventClient) ReleaseSeats(ctx context.Context, eventPublicID string, seats uint64) error {
	// TODO: Implementation
	return nil
}
