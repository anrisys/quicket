package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/anrisys/quicket/test/test_utils"
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
					"token":     "",
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
				validateGenericSuccess(t, responseBody)
				validateLoginData(t, responseBody, testUser.PublicID)
			} else {
				validateErrorResponse(t, responseBody, tt.expectedResponseBody)
			}
		})
	}
}

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
		name                 string
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
				validateGenericSuccess(t, responseBody)
			} else {
				validateErrorResponse(t, responseBody, tt.expectedResponseBody)
			}
		})
	}
}

func validateGenericSuccess(t *testing.T, response map[string]any) {
	t.Helper()

	// Check success code
	if code, exists := response["code"]; !exists || code != "SUCCESS" {
		t.Errorf("Expected code 'SUCCESS', got %v", response["code"])
	}

	// Check message
	if message, exists := response["message"]; !exists || message == "" {
		t.Error("Expected a non-empty message in success response")
	}
}

func validateErrorResponse(t *testing.T, actualResponse, expectedResponse map[string]any) {
	t.Helper()
	// Check error code exists and matches
	if actualCode, exists := actualResponse["code"]; !exists || actualCode != expectedResponse["code"] {
		t.Errorf("Expected error code %v, got %v", expectedResponse["code"], actualCode)
	}

	// Check error message exists
	if actualMessage, exists := actualResponse["message"]; !exists || actualMessage == "" {
		t.Error("Expected error message in response")
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