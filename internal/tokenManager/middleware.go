package tokenManager

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
	db_dto "main/internal/adapters/dto"
	"net/http"
	"strings"
)

func (t *tokenManager) CheckIsSeller(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tokenString string
		if val := r.Header.Get("Authorization"); val != "" {
			tokenString = strings.TrimPrefix(val, "Bearer ")
		} else {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("missing authorization header"))
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(t.signingKey), nil
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("invalid authorization token or the token has expired"))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("invalid authorization token claims"))
			return
		}

		seller, ok := claims["user"].(string)
		if !ok {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("not authorized as seller"))
			return
		}
		dto := db_dto.GetUserDTO{
			ID: seller,
		}
		user, err := t.storage.GetUser(r.Context(), dto)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("an error occurred in the database"))
			return
		}
		if !user.IsSeller {
			log.Println(user)
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("not authorized as seller"))
			return
		}

		next(w, r)
	}
}
