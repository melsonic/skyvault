package main

import (
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/melsonic/skyvault/auth/db"
	"github.com/melsonic/skyvault/auth/handler"
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
	mux.HandleFunc("POST /logout", handler.LogoutHandler)
	mux.HandleFunc("GET /authorize", handler.AuthorizeHandler)

	server := &http.Server{
		Addr:           ":8002",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(server.ListenAndServe())
}
