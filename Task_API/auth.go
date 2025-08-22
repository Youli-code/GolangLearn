package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("TASK_API_JWT_SECRET"))

func GenerateJWT(userID int) (string, error) {
	if len(jwtSecret) == 0 {
		return "", errors.New("JWT secret not set")
	}

	claims := jwt.MapClaims{
		"userId": userID,
		"exp":    time.Now().Add(2 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

type contextKey int

const userKey contextKey = iota

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			writeErr(w, http.StatusUnauthorized, "missing or invalid token")
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			writeErr(w, http.StatusUnauthorized, "invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			writeErr(w, http.StatusUnauthorized, "invalid claims")
			return
		}

		userID, ok := claims["userId"].(float64)
		if !ok {
			writeErr(w, http.StatusUnauthorized, "invalid user ID")
			return
		}

		ctx := context.WithValue(r.Context(), userKey, int(userID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserID(r *http.Request) int {
	if val, ok := r.Context().Value(userKey).(int); ok {
		return val
	}
	return 0
}
