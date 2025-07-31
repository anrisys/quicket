package security

import (
	"fmt"

	"github.com/anrisys/quicket/pkg/config"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AccountSecurityInterface interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hashedPassword string) bool
	GeneratePublicID() (string, error)
}

type AccountSecurity struct {
	bcryptCost int
}

func NewAccountSecurity(cfg *config.AppConfig) *AccountSecurity {
	return &AccountSecurity{
		bcryptCost: cfg.Security.BcryptCost,
	}
}

func (s *AccountSecurity) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.bcryptCost)
	return string(bytes), err
}

func (s *AccountSecurity) CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (s *AccountSecurity) GeneratePublicID() (string, error) {
	publicID, err := uuid.NewRandom()

	if err != nil {
		return "", fmt.Errorf("failed to generate public ID: %w", err)
	}

	return publicID.String(), nil
}