package token

import (
	"time"

	"github.com/anrisys/quicket/user-service/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

type TokenGeneratorInterface interface {
	GenerateToken(publicID, role string) (string, error)
}

type TokenGenerator struct {
	secret string
	issuer string
	expiry time.Duration
}

func NewTokenGenerator(cfg *config.Config) *TokenGenerator {
	return &TokenGenerator{
		secret: cfg.JWT.JWTSecret,
		issuer: cfg.JWT.JWTIssuer,
		expiry: cfg.JWT.JWTExpiry,
	}
}

func (g *TokenGenerator) GenerateToken(publicID, role string) (string, error) {
	claims := jwt.MapClaims{
		"sub": publicID,
		"role": role,
		"iss": g.issuer,
		"exp": time.Now().Add(g.expiry).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(g.secret))
}