package security

import (
	"context"
	"fmt"

	"github.com/anrisys/quicket/user-service/pkg/config"
	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
)

type AccountSecurityInterface interface {
	HashPassword(ctx context.Context, password string) (string, error)
	CheckPasswordHash(ctx context.Context, password, hashedPassword string) bool
	GeneratePublicID(ctx context.Context) (string, error)
}

type AccountSecurity struct {
	bcryptCost int
}

func NewAccountSecurity(cfg *config.AppConfig) *AccountSecurity {
	return &AccountSecurity{
		bcryptCost: cfg.Security.BcryptCost,
	}
}

func (s *AccountSecurity) HashPassword(_ctx context.Context, password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.bcryptCost)
	return string(bytes), err
}

func (s *AccountSecurity) CheckPasswordHash(_ctx context.Context, password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (s *AccountSecurity) GeneratePublicID(_ctx context.Context) (string, error) {
	publicID, err := uuid.NewRandom()

	if err != nil {
		return "", fmt.Errorf("failed to generate public ID: %w", err)
	}

	return publicID.String(), nil
}