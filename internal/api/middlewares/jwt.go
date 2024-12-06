package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/edulustosa/verdin/internal/api"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Middlewares struct {
	JWTKey string
}

func (m *Middlewares) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := getTokenFromHeader(r)
		if err != nil {
			unauthorized(w, err)
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}

			return []byte(m.JWTKey), nil
		})
		if err != nil {
			unauthorized(w, err)
			return
		}

		userID, err := getUserIDFromClaims(token)
		if err != nil {
			unauthorized(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), api.UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing Authorization header")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		return "", errors.New("missing token")
	}

	return token, nil
}

func unauthorized(w http.ResponseWriter, err error) {
	api.Encode(w, http.StatusUnauthorized, api.Errors{
		Errors: []api.Error{
			{
				StatusCode: http.StatusUnauthorized,
				Message:    err.Error(),
			},
		},
	})
}

func getUserIDFromClaims(token *jwt.Token) (uuid.UUID, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, err := uuid.Parse(claims["sub"].(string))
		if err != nil {
			return uuid.Nil, errors.New("invalid user id in claims")
		}

		return userID, nil
	}

	return uuid.Nil, errors.New("invalid token claims")
}
