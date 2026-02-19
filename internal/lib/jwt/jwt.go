package jwt

import (
	"time"

	"github.com/Kefir4c/sso-service/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
)

func NewToken(user *models.User, app *models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["user_id"] = user.ID
	claims["email"] = user.Email
	claims["app_id"] = app.ID
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SigningString()
	if err != nil {
		return "", err
	}
	return tokenString, nil

}
