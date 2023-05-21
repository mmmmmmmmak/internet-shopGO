package tokenManager

import (
	"github.com/golang-jwt/jwt"
	"time"
)

type tokenManager struct {
	signingKey string
}

func (t *tokenManager) GenerateAccessToken(userID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(1 * time.Hour).Unix()
	claims["authorized"] = true
	claims["user"] = userID

	tokenString, err := token.SignedString([]byte(t.signingKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (t *tokenManager) GenerateRefreshToken() (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()

	tokenString, err := token.SignedString([]byte(t.signingKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func NewTokenManager(secret string) *tokenManager {
	return &tokenManager{signingKey: secret}
}
