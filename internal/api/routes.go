package api

import (
	"net/http"
	"tt/internal/handlers"
	"tt/internal/middleware"
	"tt/internal/services"
)

func SetupRouter(authService services.AuthServiceInterface, handlers handlers.Handlers) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("POST /api/register", handlers.Auth.Register)
	router.HandleFunc("POST /api/login", handlers.Auth.Login)
	router.HandleFunc("POST /api/logout", middleware.AuthMiddleware(authService, handlers.Auth.Logout))

	router.HandleFunc("POST /api/tables", middleware.AuthMiddleware(authService, handlers.DynamicData.CreateTable))
	router.HandleFunc("POST /api/tables/{tableID}/rows", middleware.AuthMiddleware(authService, handlers.DynamicData.CreateRow))
	router.HandleFunc("GET /api/tables/{tableID}/rows", middleware.AuthMiddleware(authService, handlers.DynamicData.GetRows))
	router.HandleFunc("PUT /api/tables/{tableID}/rows/{rowID}", middleware.AuthMiddleware(authService, handlers.DynamicData.UpdateRow))
	router.HandleFunc("DELETE /api/tables/{tableID}/rows/{rowID}", middleware.AuthMiddleware(authService, handlers.DynamicData.DeleteRow))

	return router
}
