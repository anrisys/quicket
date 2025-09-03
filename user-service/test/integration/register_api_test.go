package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/anrisys/quicket/user-service/test/test_utils"
)

func TestRegister(t *testing.T) {
	s := test_utils.NewTestServer()
	defer s.Close()

	// Seed user for duplicate case
	_, err := test_utils.CreateTestUser(s.App.Config)
	if err != nil {
		t.Fatalf("failed to seed user data %v", err)
	}

	defer test_utils.CleanupTestUser(s.App.Config, "newuser@example.com")
	defer test_utils.CleanupTestUser(s.App.Config, "user@example.com")

	tests := []struct {
		name     string
		payload              map[string]string
		expectedStatus       int
		expectedResponseBody map[string]any
		expectError          bool
	}{
		{
			name: "register success",
			payload: map[string]string{
				"email": "newuser@example.com",
				"password": "Pass345!@#",
				"password_confirmation": "Pass345!@#",
				"role": "user",
			},
			expectedStatus: http.StatusCreated,
			expectedResponseBody: map[string]any{
				"code": "SUCCESS",
				"message": "User registered successful",
			},
			expectError: false,
		},
		{
			name: "register failed - duplicate email",
			payload: map[string]string{
				"email": "user@example.com",
				"password": "Pass345!@#",
				"password_confirmation": "Pass345!@#",
				"role": "user",
			},
			expectedStatus: http.StatusConflict,
			expectedResponseBody: map[string]any{
				"code": "CONFLICT_ERROR",
				"message": "email already registered",
			},
			expectError: true,
		},
		{
			name: "register failed - invalid email format",
			payload: map[string]string{
				"email": "not-an-email",
				"password": "Pass345!@#",
				"password_confirmation": "Pass345!@#",
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponseBody: map[string]any{
				"code": "VALIDATION_ERROR",
				"message": "Invalid login data",
			},
			expectError: true,
		},
		{
			name: "register failed - password mismatch",
			payload: map[string]string{
				"email": "mismatch@example.com",
				"password": "Pass345!@#",
				"password_confirmation": "Different123!",
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponseBody: map[string]any{
				"code": "VALIDATION_ERROR",
				"message": "Invalid login data",
			},
			expectError: true,
		},
		{
			name: "register failed - weak password",
			payload: map[string]string{
				"email": "weakpass@example.com",
				"password": "123",
				"password_confirmation": "123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponseBody: map[string]any{
				"code": "VALIDATION_ERROR",
				"message": "Invalid login data",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := s.MakeRequest("POST", "/api/v1/register", "", tt.payload)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			var responseBody map[string]any
			if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if !tt.expectError {
				test_utils.ValidateGenericSuccess(t, responseBody)
			} else {
				test_utils.ValidateErrorResponse(t, responseBody, tt.expectedResponseBody)
			}
		})
	}
}