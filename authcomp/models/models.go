package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Email               string    `json:"email"`
	Password            string    `json:"password"`
	Name                string    `json:"name"`
	Gender              string    `json:"gender"`
	DateOfBirth         time.Time `json:"date_of_birth"`
	DateCreated         time.Time `json:"date_created"`
	RefreshTokenVersion int       `json:"refresh_token_version"`
}

type AuthResponse struct {
	AccessToken  string `json:"accesstoken"`
	RefreshToken string `json:"refreshtoken"`
}

type JWTTokenClaims struct {
	Email               string
	TokenType           string
	RefreshTokenVersion int
	jwt.RegisteredClaims
}
