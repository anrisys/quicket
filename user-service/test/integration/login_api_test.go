package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/anrisys/quicket/user-service/test/test_utils"
)

func TestLogin(t *testing.T) {
	s := test_utils.NewTestServer()
	defer s.Close()

	// Seed user test data
	testUser, err := test_utils.CreateTestUser(s.App.Config)
	if err != nil {
		t.Fatalf("failed to seed user data %v", err)
	}

	defer test_utils.CleanupTestUser(s.App.Config, testUser.Email)

	tests := []struct {
		name                 string
		payload              map[string]string
		expectedStatus       int
		expectedResponseBody map[string]any
		expectError          bool
	}{
		{
			name: "login success",
			payload: map[string]string{
				"email": "user@example.com",
				"password": "Pass345!@#",
			},
			expectedStatus: http.StatusOK,
			expectedResponseBody: map[string]any{
				"code": "SUCCESS",
				"message": "User logged in successful",
				"data": map[string]any{
					"public_id": testUser.PublicID,
					"token":  "",
				},
			},
			expectError: false,
		},
		{
			name: "login failed - wrong password",
			payload: map[string]string{
				"email": "user@example.com",
				"password": "wrongpassword",
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponseBody: map[string]any{
				"code": "VALIDATION_ERROR",
				"message": "email or password is wrong",
			},
			expectError: true,
		},
		{
			name: "login failed - user not found",
			payload: map[string]string{
				"email": "nonexistent@example.com",
				"password": "anypassword",
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponseBody: map[string]any{
				"code": "VALIDATION_ERROR",
				"message": "email or password is wrong",
			},
			expectError: true,
		},
		{
			name: "login failed - invalid email format",
			payload: map[string]string{
				"email": "invalid-email",
				"password": "Pass345!@#",
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
			resp, err := s.MakeRequest("POST", "/api/v1/login", "", tt.payload)
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
				validateLoginData(t, responseBody, testUser.PublicID)
			} else {
				test_utils.ValidateErrorResponse(t, responseBody, tt.expectedResponseBody)
			}
		})
	}
}

func validateLoginData(t *testing.T, response map[string]any, expectedPublicID string) {
	t.Helper()

	// Check data structure
	data, dataExists := response["data"].(map[string]any)
	if !dataExists {
		t.Error("Expected 'data' field in success response")
		return
	}

	// Check token
	if token, tokenExists := data["token"]; !tokenExists || token == "" {
		t.Error("Expected a non-empty 'token' in response data")
	}

	// Check public ID
	if publicID, idExists := data["public_id"]; !idExists || publicID != expectedPublicID {
		t.Errorf("Expected 'public_id' %s, got %v", expectedPublicID, publicID)
	}
}