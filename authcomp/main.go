package main

import (
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/melsonic/skyvault/auth/db"
	"github.com/melsonic/skyvault/auth/handler"
	"github.com/melsonic/skyvault/auth/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = db.InitDB()
	if err != nil {
		log.Fatal("Error in InitDB")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth/register	", handler.RegisterHandler)
	mux.HandleFunc("POST /auth/login", handler.LoginHandler)
	mux.HandleFunc("POST /auth/logout", middleware.AuthMiddleware(handler.LogoutHandler))
	mux.HandleFunc("POST /auth/refresh", middleware.AuthMiddleware(handler.NewAccessTokenHandler))
	mux.HandleFunc("GET /users/whoami", middleware.AuthMiddleware(handler.WhoamiHandler))
	mux.HandleFunc("GET /user/me", middleware.AuthMiddleware(handler.GetUserHandler))
	mux.HandleFunc("PUT /user/me", middleware.AuthMiddleware(handler.UpdateUserHandler))
	mux.HandleFunc("DELETE /user/me", middleware.AuthMiddleware(handler.DeleteUserHandler))
	mux.HandleFunc("POST /auth/password-reset", handler.PasswordResetHandler)
	mux.HandleFunc("GET /auth/password-reset/{hash}", handler.UpdatePasswordFormHandler)
	mux.HandleFunc("POST /auth/password-reset/{hash}", handler.UpdatePasswordHandler)

	server := &http.Server{
		Addr:           ":8003",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(server.ListenAndServe())
}
