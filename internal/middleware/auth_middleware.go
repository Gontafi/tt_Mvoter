package middleware

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strconv"
	"strings"
	"tt/internal/services"
	"tt/pkg/utils"
)

func AuthMiddleware(authService services.AuthServiceInterface, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			utils.SendError(w, "missing token", http.StatusBadRequest)
			return
		}

		bearer := strings.Split(tokenString, " ")
		if len(bearer) < 2 {
			utils.SendError(w, "missing token or user_id", http.StatusBadRequest)
			return
		}

		token, err := jwt.Parse(bearer[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("HMAC_SECRET")), nil
		})

		if err != nil {
			utils.SendError(w, "Your Token is invalid or expired", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			utils.SendError(w, "Invalid Token Claims", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["uid"].(float64)

		if !ok {
			utils.SendError(w, "Failed to get userID", http.StatusUnauthorized)
			return
		}

		redisToken, err := authService.GetTokenFromRedis(r.Context(), bearer[1], strconv.FormatUint(uint64(userID), 10))
		if err != nil || redisToken != bearer[1] {
			utils.SendError(w, "Token not found or expired in Redis", http.StatusUnauthorized)
			return
		}

		handler.ServeHTTP(w, r)
	}
}
