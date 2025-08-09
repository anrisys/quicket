package token

import (
	"time"

	"github.com/anrisys/quicket/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

type Generator struct {
	secret string
	issuer string
	expiry time.Duration
}

type Claims struct {
	sub 	string
	role	string
	iss 	string
	iat 	string
}

func NewGenerator(cfg *config.AppConfig) *Generator {
	return &Generator{
		secret: cfg.Security.JWTSecret,
		issuer: cfg.Security.JWTIssuer,
		expiry: cfg.Security.JWTExpiry,
	}
}

func (g *Generator) GenerateToken(publicID, role string) (string, error) {
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