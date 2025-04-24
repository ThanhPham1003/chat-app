package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(username, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString([]byte(secret))
}
