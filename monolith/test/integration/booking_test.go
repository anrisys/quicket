package integration_test

import (
	"encoding/json"
	"net/http"
	"testing"

	tokenGenerator "github.com/anrisys/quicket/pkg/token"
	"github.com/anrisys/quicket/test/test_utils"
)

func TestCreateBooking(t *testing.T) {
	s := test_utils.NewTestServer()
	defer s.Close()

	testUser, err := test_utils.CreateTestUser(s.App.Config)
	if err != nil {
		t.Logf("failed to seed user data %v", err)
	}

	tg := tokenGenerator.NewGenerator(s.App.Config)
	tokenTest, err := tg.GenerateToken(testUser.PublicID, testUser.Role)
	if err != nil {
		t.Logf("failed to generate test user token %v", err)
	}
	validEvent, err := test_utils.CreateTestEvent(s.App.Config, "Go Tech Conference 2025", testUser.ID)
	if err != nil {
		t.Logf("failed to seed event data %v", err)
	}
	limitedSeatsEv, err := test_utils.CreateLimitedSeatsEventTest(s.App.Config, "Limited Seats Event", testUser.ID)
	if err != nil {
		t.Logf("failed to seed limited time seats event %v", err)
	}
	pastEv, err := test_utils.CreatePastEventTest(s.App.Config, "Past event", testUser.ID)
	if err != nil {
		t.Logf("failed to seed past even test %v", err)
	}
	defer test_utils.CleanupTestUser(s.App.Config, testUser.Email)
	defer test_utils.CleanupTestEvent(s.App.Config, validEvent.Title)
	defer test_utils.CleanupTestEvent(s.App.Config, pastEv.Title)
	defer test_utils.CleanupTestEvent(s.App.Config, limitedSeatsEv.Title)

	tests := []struct{
		name                 string
		token                string
		eventID              string
		payload              map[string]any
		expectedStatus       int
		expectedResponseBody map[string]any
		expectError          bool
	} {
		{
			name: "create booking success",
			token: tokenTest,
			eventID: validEvent.PublicID,
			payload: map[string]any{
				"seats": uint(2),
			},
			expectedStatus: http.StatusCreated,
			expectedResponseBody: map[string]any{
				"code": "SUCCESS",
				"message": "Booking created successfully",
				"booking": map[string]any{
					"event_id": validEvent.PublicID,
					"user_id": testUser.PublicID,
					"seats": float64(2),
				},
			},
			expectError: false,
		},
		{
			name: "create booking failed - past event",
			token: tokenTest,
			eventID: pastEv.PublicID,
			payload: map[string]any{
				"seats": uint(1),
			},
			expectedStatus: http.StatusConflict,
			expectedResponseBody: map[string]any{
				"code": "CONFLICT_ERROR",
				"message": "can not book past event",
			},
			expectError: true,
		},
		{
			name: "create booking failed - not enough seats",
			token: tokenTest,
			eventID: limitedSeatsEv.PublicID,
			payload: map[string]any{
				"seats": uint(6),
			},
			expectedStatus: http.StatusConflict,
			expectedResponseBody: map[string]any{
				"code": "CONFLICT_ERROR",
				"message": "not enough available seats",
			},
			expectError: true,
		},
		{
			name: "create booking failed - non-existent event",
			token: tokenTest,
			eventID: "non-existent-id",
			payload: map[string]any{
				"seats": uint(1),
			},
			expectedStatus: http.StatusNotFound,
			expectedResponseBody: map[string]any{
				"code": "NOT_FOUND",
				"message": "event not found",
			},
			expectError: true,
		},
		{
			name: "create booking failed - no seats",
			token: tokenTest,
			eventID: validEvent.PublicID,
			payload: map[string]any{
				"seats": uint(0),
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponseBody: map[string]any{
				"code": "VALIDATION_ERROR",
				"message": "Invalid login data",
			},
			expectError: true,
		},
		{
			name: "create booking failed - invalid payload type",
			token: tokenTest,
			eventID: validEvent.PublicID,
			payload: map[string]any{
				"seats": "not a number",
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponseBody: map[string]any{
				"code": "VALIDATION_ERROR",
				"message": "Invalid login data",
			},
			expectError: true,
		},
		{
			name: "create booking failed - unauthorized",
			token: "",
			eventID: validEvent.PublicID,
			payload: map[string]any{
				"seats": uint(1),
			},
			expectedStatus: http.StatusUnauthorized,
			expectedResponseBody: map[string]any{
				"code": "UNAUTHORIZED",
				"message": "Authentication required",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := s.MakeRequest("POST", "/api/v1/bookings/"+tt.eventID, tt.token, tt.payload)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Test: %s", tt.expectedStatus, resp.StatusCode, tt.name)
			}

			var responseBody map[string]any
			if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}
			t.Logf("==RESPONSE BODY JSON== %s", responseBody)
			if !tt.expectError {
				validateBookingSuccessCreate(t, responseBody, tt.expectedResponseBody)
			} else {
				test_utils.ValidateErrorResponse(t, responseBody, tt.expectedResponseBody)
			}
		})
	}
}

func validateBookingSuccessCreate(t *testing.T, actualResponse, expectedResponseBody map[string]any) {
	// First, check for the overall success structure
	test_utils.ValidateGenericSuccess(t, actualResponse)
	
	// Assert the booking object
	actualBooking, ok := actualResponse["booking"].(map[string]any)
	if !ok {
		t.Fatalf("expected 'booking' field to be a map, got %T", actualResponse["booking"])
	}
	
	// Assert booking fields
	expectedBooking, ok := expectedResponseBody["booking"].(map[string]any)
	if !ok {
		t.Fatalf("expectedResponseBody does not contain a 'booking' map")
	}

	if actualBooking["event_id"] != expectedBooking["event_id"] {
		t.Errorf("expected event_id %v, got %v", expectedBooking["event_id"], actualBooking["event_id"])
	}

	if actualBooking["user_id"] != expectedBooking["user_id"] {
		t.Errorf("expected user_id %v, got %v", expectedBooking["user_id"], actualBooking["user_id"])
	}

	// JSON numbers are unmarshaled as float64, so we need to compare float64 values
	if actualBooking["seats"] != expectedBooking["seats"] {
		t.Errorf("expected number of seats %v, got %v", expectedBooking["seats"], actualBooking["seats"])
	}
	
	// Validate dynamic fields
	publicID, ok := actualBooking["id"].(string)
	if !ok || publicID == "" {
		t.Errorf("expected booking 'id' to be a non-empty string, got %v", actualBooking["id"])
	}
	
	status, ok := actualBooking["status"].(string)
	if !ok || status == "" {
		t.Errorf("expected booking 'status' to be a non-empty string, got %v", actualBooking["status"])
	}

	expiredAt, ok := actualBooking["expired_at"].(string)
	if !ok || expiredAt == "" {
		t.Errorf("expected booking 'expired_at' to be a non-empty string, got %v", actualBooking["expired_at"])
	}
}