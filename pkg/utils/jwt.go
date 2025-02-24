package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
	"time"
)

func GenerateJWT(id uint64) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["uid"] = id
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString([]byte(os.Getenv("HMAC_SECRET")))

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func ParseUserIDJWTInHandler(w http.ResponseWriter, r *http.Request) int64 {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		SendError(w, "missing token", http.StatusBadRequest)
		return 0
	}

	bearer := strings.Split(tokenString, " ")
	if len(bearer) < 2 {
		SendError(w, "missing token or user_id", http.StatusBadRequest)
		return 0
	}

	token, err := jwt.Parse(bearer[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("HMAC_SECRET")), nil
	})

	if err != nil {
		SendError(w, "Your Token is invalid or expired", http.StatusUnauthorized)
		return 0
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		SendError(w, "Invalid Token Claims", http.StatusUnauthorized)
		return 0
	}

	userID, ok := claims["uid"].(float64)

	if !ok {
		SendError(w, "Failed to get userID", http.StatusUnauthorized)
		return 0
	}

	return int64(userID)
}
