package token

import (
	"errors"
	"time"
	"user-management-api/internal/core/port"

	"github.com/golang-jwt/jwt/v5"
)

type jwtProvider struct {
	secretKey []byte
	ttl       time.Duration
}

func NewJWTProvider(secretKey string, ttl time.Duration) port.TokenProvider {
	return &jwtProvider{secretKey: []byte(secretKey), ttl: ttl}
}

func (p *jwtProvider) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(p.ttl).Unix(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(p.secretKey)
}

func (p *jwtProvider) ValidateToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return p.secretKey, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("invalid token subject")
	}

	return userID, nil
}
