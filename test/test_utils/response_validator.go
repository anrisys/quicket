package test_utils

import "testing"

func ValidateGenericSuccess(t *testing.T, response map[string]any) {
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

func ValidateErrorResponse(t *testing.T, actualResponse, expectedResponse map[string]any) {
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