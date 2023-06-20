package tokenManager

import (
	"context"
	"github.com/golang-jwt/jwt"
	db_dto "main/internal/adapters/dto"
	"main/internal/domain/entity"
	"time"
)

type UserStorage interface {
	GetUser(ctx context.Context, user db_dto.GetUserDTO) (entity.User, error)
}

type tokenManager struct {
	signingKey string
	storage    UserStorage
}

func (t *tokenManager) GenerateAccessToken(userID string, isSeller bool) (string, error) {
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

func NewTokenManager(secret string, storage UserStorage) *tokenManager {
	return &tokenManager{
		signingKey: secret,
		storage:    storage,
	}
}
