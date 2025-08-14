package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Email               string    `json:"email,omitempty"`
	Password            string    `json:"password,omitempty"`
	Name                string    `json:"name,omitempty"`
	Gender              string    `json:"gender,omitempty"`
	DateOfBirth         time.Time `json:"date_of_birth,omitempty"`
	DateCreated         time.Time `json:"date_created,omitempty"`
	RefreshTokenVersion int       `json:"refresh_token_version,omitempty"`
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
	ID            string    `json:"id,omitempty"`
	Email         string    `json:"email,omitempty"`
	RequestExpiry time.Time `json:"request_expiry,omitempty"`
	UsedFlag      bool      `json:"used,omitempty"`
}

type EmailTemplateData struct {
	ResetURL string
}

type PasswordResetUpdateRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type ResponseMessage struct {
	Message string `json:"message"`
}
