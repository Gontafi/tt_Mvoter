package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"tt/internal/models"
	"tt/internal/services"
)

func sendError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	var err models.Error
	err.SetError(message)
	json.NewEncoder(w).Encode(err)
}

func AuthMiddleware(authService *services.AuthService, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Token")
		if tokenString == "" {
			sendError(w, "No Token Found", http.StatusUnauthorized)
			return
		}

		mySigningKey := []byte(os.Getenv("HMAC_SECRET"))

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return mySigningKey, nil
		})

		if err != nil {
			sendError(w, "Your Token is invalid or expired", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			sendError(w, "Invalid Token Claims", http.StatusUnauthorized)
			return
		}
		userID, ok := claims["uid"].(string)

		if !ok {
			sendError(w, "Failed to get userID", http.StatusUnauthorized)
			return
		}

		redisToken, err := authService.GetTokenFromRedis(r.Context(), tokenString, userID)
		if err != nil || redisToken != tokenString {
			sendError(w, "Token not found or expired in Redis", http.StatusUnauthorized)
			return
		}

		handler.ServeHTTP(w, r)
	}
}
