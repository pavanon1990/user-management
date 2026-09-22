package router

import (
	"net/http"
	"user-management-api/internal/adapter/handler"
	"user-management-api/internal/adapter/middleware"
	"user-management-api/internal/core/port"
)

func NewRouter(userHandler *handler.UserHandler, tokenProvider port.TokenProvider) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /auth/register", userHandler.RegisterUser)
	mux.HandleFunc("POST /auth/login", userHandler.LoginUser)

	auth := middleware.Auth(tokenProvider)

	mux.HandleFunc("GET /users/{id}", auth(userHandler.GetUserByID))
	mux.HandleFunc("GET /users", auth(userHandler.GetUsers))
	mux.HandleFunc("PATCH /users/{id}", auth(userHandler.UpdateUser))
	mux.HandleFunc("DELETE /users/{id}", auth(userHandler.DeleteUser))

	return middleware.Log(mux)
}
