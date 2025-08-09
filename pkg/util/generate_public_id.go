package util

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func GeneratePublicID(_ctx context.Context) (string, error) {
	publicID, err := uuid.NewRandom()

	if err != nil {
		return "", fmt.Errorf("failed to generate public ID: %w", err)
	}

	return publicID.String(), nil
}

