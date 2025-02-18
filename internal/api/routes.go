package api

import (
	"net/http"
	"tt/internal/handlers"
	"tt/internal/middleware"
	"tt/internal/services"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Allow frontend
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Preflight Request Handling
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func SetupRouter(authService services.AuthServiceInterface, handlers handlers.Handlers) http.Handler {
	router := http.NewServeMux()

	router.Handle("POST /api/register", http.HandlerFunc(handlers.Auth.Register))
	router.Handle("POST /api/login", http.HandlerFunc(handlers.Auth.Login))
	router.Handle("POST /api/logout", middleware.AuthMiddleware(authService, http.HandlerFunc(handlers.Auth.Logout)))

	router.Handle("POST /api/tables", middleware.AuthMiddleware(authService, http.HandlerFunc(handlers.DynamicData.CreateTable)))
	router.Handle("POST /api/tables/{tableID}/rows", middleware.AuthMiddleware(authService, http.HandlerFunc(handlers.DynamicData.CreateRow)))
	router.Handle("GET /api/tables/{tableID}/rows", middleware.AuthMiddleware(authService, http.HandlerFunc(handlers.DynamicData.GetRows)))
	router.Handle("PUT /api/tables/{tableID}/rows/{rowID}", middleware.AuthMiddleware(authService, http.HandlerFunc(handlers.DynamicData.UpdateRow)))
	router.Handle("DELETE /api/tables/{tableID}/rows/{rowID}", middleware.AuthMiddleware(authService, http.HandlerFunc(handlers.DynamicData.DeleteRow)))

	// Wrap the entire router with CORS middleware
	return corsMiddleware(router)
}
