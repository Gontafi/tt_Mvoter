package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strconv"
	"strings"
	"tt/pkg/utils"

	"tt/internal/models"
	"tt/internal/services"
)

type AuthHandler struct {
	service services.AuthServiceInterface
}

func NewAuthHandler(service services.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.SendError(w, "invalid input", http.StatusBadRequest)
		return
	}

	_, err := h.service.Register(r.Context(), user)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.SendError(w, "invalid input", http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(r.Context(), user)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": *token})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		utils.SendError(w, "missing token or user_id", http.StatusBadRequest)
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
		utils.SendError(w, "Invalid Token Claims uid", http.StatusBadRequest)
	}

	if err := h.service.Logout(r.Context(), bearer[1], strconv.FormatUint(uint64(userID), 10)); err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
