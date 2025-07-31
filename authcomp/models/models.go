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
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type JWTTokenClaims struct {
	Email               string
	Name                string
	TokenType           string
	RefreshTokenVersion int
	jwt.RegisteredClaims
}

type NewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type NewAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type PasswordResetRequest struct {
	ID          string    `json:"id,omitempty"`
	Email       string    `json:"email,omitempty"`
	RequestedAt time.Time `json:"requested_at,omitempty"`
	UsedFlag    bool      `json:"used,omitempty"`
}

type EmailTemplateData struct {
	ResetURL string
}