package integration

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	tokenGenerator "github.com/anrisys/quicket/pkg/token"
	"github.com/anrisys/quicket/test/test_utils"
)

func TestCreateEvent(t *testing.T) {
	s := test_utils.NewTestServer()
	defer s.Close()

	testUser, err := test_utils.CreateTestUser(s.App.Config)
	if err != nil {
		t.Fatalf("failed to seed user data %v", err)
	}
	tg := tokenGenerator.NewGenerator(s.App.Config)
	tokenTest, err := tg.GenerateToken(testUser.PublicID, testUser.Role)
	if err != nil {
		t.Fatalf("failed to generate test user token %v", err)
	}
	duplicateTitle := "Duplicate ttile"
	ev, err := test_utils.CreateTestEvent(s.App.Config, duplicateTitle, testUser.ID)
	if err != nil {
		t.Fatalf("failed to seed event data %v", err)
	}
	defer test_utils.CleanupTestUser(s.App.Config, testUser.Email)
	defer test_utils.CleanupTestEvent(s.App.Config, ev.Title)
	tomorrow := time.Now().Add(24 * time.Hour)
	dayAfterTomorrow := tomorrow.Add(24 * time.Hour)

	tests := []struct {
		name                 string
		token               string
		payload             map[string]any
		expectedStatus      int
		expectedResponse    map[string]any
		expectError         bool
		description         string
	}{
		// SUCCESS CASES
		{
			name: "create event success",
			token: tokenTest,
			payload: map[string]any{
				"title": "Tech Conference 2024",
				"start_date": tomorrow.Format(time.RFC3339),
				"end_date": dayAfterTomorrow.Format(time.RFC3339),
				"description": "Annual technology conference",
				"max_seats": 100,
			},
			expectedStatus: http.StatusCreated,
			expectedResponse: map[string]any{
				"code":    "SUCCESS",
				"message": "Event created successfully",
				"public_id": "", // Will check if not empty
				"title":   "Tech Conference 2024",
			},
			expectError: false,
			description: "Valid event creation should succeed",
		},
		{
			name:   "create event without description",
			token:  tokenTest,
			payload: map[string]any{
				"title":      "Music Concert",
				"start_date": tomorrow.Format(time.RFC3339),
				"end_date":   dayAfterTomorrow.Format(time.RFC3339),
				"max_seats":  500,
			},
			expectedStatus: http.StatusCreated,
			expectedResponse: map[string]any{
				"code":    "SUCCESS",
				"message": "Event created successfully",
			},
			expectError: false,
			description: "Event without description should succeed",
		},
		
		// AUTHENTICATION FAILURES
		{
			name:   "create event unauthorized - no token",
			token:  "",
			payload: map[string]any{
				"title":      "Unauthorized Event",
				"start_date": tomorrow.Format(time.RFC3339),
				"end_date":   dayAfterTomorrow.Format(time.RFC3339),
				"max_seats":  50,
			},
			expectedStatus: http.StatusUnauthorized,
			expectedResponse: map[string]any{
				"code":    "UNAUTHORIZED",
				"message": "Authentication required",
			},
			expectError: true,
			description: "Should reject request without authentication token",
		},
		{
			name:   "create event unauthorized - invalid token",
			token:  "invalid-token-here",
			payload: map[string]any{
				"title":      "Invalid Token Event",
				"start_date": tomorrow.Format(time.RFC3339),
				"end_date":   dayAfterTomorrow.Format(time.RFC3339),
				"max_seats":  50,
			},
			expectedStatus: http.StatusUnauthorized,
			expectedResponse: map[string]any{
				"code":    "UNAUTHORIZED",
				"message": "Invalid token",
			},
			expectError: true,
			description: "Should reject request with invalid token",
		},
		// VALIDATION FAILURES
		{
			name:   "create event failed - title too short",
			token:  tokenTest,
			payload: map[string]any{
				"title":      "ab", // Too short
				"start_date": tomorrow.Format(time.RFC3339),
				"end_date":   dayAfterTomorrow.Format(time.RFC3339),
				"max_seats":  100,
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: map[string]any{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid request data",
			},
			expectError: true,
			description: "Should reject event with title less than 3 characters",
		},
		{
			name:   "create event failed - start date in past",
			token:  tokenTest,
			payload: map[string]any{
				"title":      "Past Event",
				"start_date": time.Now().Add(-24 * time.Hour).Format(time.RFC3339), // Yesterday
				"end_date":   tomorrow.Format(time.RFC3339),
				"max_seats":  100,
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: map[string]any{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid request data",
			},
			expectError: true,
			description: "Should reject event with start date in the past",
		},
		{
			name:   "create event failed - end date before start date",
			token:  tokenTest,
			payload: map[string]any{
				"title":      "Invalid Dates Event",
				"start_date": dayAfterTomorrow.Format(time.RFC3339),
				"end_date":   tomorrow.Format(time.RFC3339), // End before start
				"max_seats":  100,
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: map[string]any{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid request data",
			},
			expectError: true,
			description: "Should reject event with end date before start date",
		},
		{
			name:   "create event failed - zero seats",
			token:  tokenTest,
			payload: map[string]any{
				"title":      "Zero Seats Event",
				"start_date": tomorrow.Format(time.RFC3339),
				"end_date":   dayAfterTomorrow.Format(time.RFC3339),
				"max_seats":  0, // Invalid
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: map[string]any{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid request data",
			},
			expectError: true,
			description: "Should reject event with zero max seats",
		},
		// DUPLICATE TITLE
		{
			name:   "create event failed - duplicate title",
			token:  tokenTest,
			payload: map[string]any{
				"title": duplicateTitle,
				"start_date": tomorrow.Add(48 * time.Hour).Format(time.RFC3339),
				"end_date": dayAfterTomorrow.Add(48 * time.Hour).Format(time.RFC3339),
				"max_seats": 200,
			},
			expectedStatus: http.StatusConflict,
			expectedResponse: map[string]any{
				"code":    "CONFLICT_ERROR",
				"message": "event with this title already exists",
			},
			expectError: true,
			description: "Should reject event with duplicate title",
		},
	}
		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := s.MakeRequest("POST", "/api/v1/events", tt.token, tt.payload)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()
			
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Test: %s", 
					tt.expectedStatus, resp.StatusCode, tt.description)
			}

			var responseBody map[string]any
			if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}
			if !tt.expectError {
				test_utils.ValidateGenericSuccess(t, responseBody)
				validateEventSuccessData(t, responseBody, tt.expectedResponse)
			} else {
				test_utils.ValidateErrorResponse(t, responseBody, tt.expectedResponse)
			}
		})
	}
}

func validateEventSuccessData(t *testing.T, actualResponse, expectedResponse map[string]any) {
	// Check event data exists
	eventData, exists := actualResponse["event"].(map[string]any)
	if !exists {
		t.Error("Expected event data in success response")
		return
	}

	// Check required event fields
	if publicID, exists := eventData["public_id"]; !exists || publicID == "" {
		t.Error("Expected non-empty public_id in event response")
	}

	if title, exists := eventData["title"]; !exists || title == "" {
		t.Error("Expected non-empty title in event response")
	}

	// Verify specific expected values if provided
	if expectedTitle, exists := expectedResponse["title"]; exists {
		if actualTitle := eventData["title"]; actualTitle != expectedTitle {
			t.Errorf("Expected title %v, got %v", expectedTitle, actualTitle)
		}
	}
}
