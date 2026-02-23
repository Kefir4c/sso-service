package jwt

import (
	"fmt"
	"time"

	"github.com/Kefir4c/sso-service/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	AppID  int    `json:"app_id"`
	jwt.RegisteredClaims
}

func NewToken(user *models.User, app *models.App, duration time.Duration) (string, error) {
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		AppID:  app.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(app.Secret))
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte("dummy"), nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)

	if !ok {
		return nil, fmt.Errorf("invalid claims type")
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func ValidateTokenWithSecret(tokenString, secret string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("validate token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid claims type")
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
