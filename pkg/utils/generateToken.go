package utils

import (
	"github.com/golang-jwt/jwt"
	"time"
)

func GenerateToken(userID string, secret []byte) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(1 * time.Hour)
	claims["authorized"] = true
	claims["user"] = userID

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
