package tokenManager

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
	"time"
)

func (t *tokenManager) TokenExpires(accessToken string) (int64, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(t.signingKey), nil
	})
	if err != nil {
		return time.Now().Unix(), err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return time.Now().Unix(), fmt.Errorf("error get user claims from token")
	}
	log.Println(claims["exp"])
	return int64(claims["exp"].(float64)), nil
}

func (t *tokenManager) TokenUser(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(t.signingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("error get user claims from token")
	}

	return claims["user"].(string), nil
}

func (t *tokenManager) ValidateToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(t.signingKey), nil
	})
	if err != nil {
		return false, fmt.Errorf("failed to parse token: %v", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return false, fmt.Errorf("invalid token claims")
	}
	exp, ok := claims["exp"].(float64)
	if !ok {
		return false, fmt.Errorf("invalid expiration time")
	}
	if int64(exp) < time.Now().Unix() {
		return false, fmt.Errorf("token is expired")
	}
	return true, nil
}
