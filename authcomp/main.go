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
	mux.HandleFunc("POST /signup", handler.SignUpHandler)
	mux.HandleFunc("POST /login", handler.LoginHandler)
	mux.HandleFunc("POST /logout", middleware.AuthMiddleware(handler.LogoutHandler))
	mux.HandleFunc("GET /authorize", middleware.AuthMiddleware(handler.AuthorizeHandler))
	mux.HandleFunc("POST /token", middleware.AuthMiddleware(handler.NewTokenHandler))
	mux.HandleFunc("GET /user/{userid}", middleware.AuthMiddleware(handler.GetUserHandler))
	mux.HandleFunc("POST /user/{userid}", middleware.AuthMiddleware(handler.UpdateUserHandler))
	mux.HandleFunc("DELETE /user/{userid}", middleware.AuthMiddleware(handler.DeleteUserHandler))
	mux.HandleFunc("POST /user/passwordreset", handler.ResetPasswordHandler)
	mux.HandleFunc("POST /user/updatepassword/{hashid}", handler.UpdatePasswordHandler)

	server := &http.Server{
		Addr:           ":8003",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(server.ListenAndServe())
}
