package jwtauth

import (
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/melsonic/skyvault/auth/db"
	"github.com/melsonic/skyvault/auth/models"
)

const (
	accessToken  = "access-token"
	refreshToken = "refresh-token"
)

func GenerateAccessToken(email string) (string, error) {
	signingKey := []byte(os.Getenv("SECRET_SIGNATURE"))
	issuer := os.Getenv("TOKEN_ISSUER")

	// create access token
	accessTokenBase := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		models.JWTTokenClaims{
			Email:     email,
			TokenType: accessToken,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(5 * time.Minute)},
				Issuer:    issuer,
			},
		},
	)
	accessToken, err := accessTokenBase.SignedString(signingKey)

	if err != nil {
		slog.Error("error creating access token", "error", err.Error())
		return "", err
	}

	return accessToken, nil
}

func GenerateRefreshToken(email string, refreshTokenVersion int) (string, error) {
	signingKey := []byte(os.Getenv("SECRET_SIGNATURE"))
	issuer := os.Getenv("TOKEN_ISSUER")

	refreshTokenBase := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		models.JWTTokenClaims{
			Email:               email,
			TokenType:           refreshToken,
			RefreshTokenVersion: refreshTokenVersion,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(5 * 24 * time.Hour)},
				Issuer:    issuer,
			},
		},
	)
	refreshToken, err := refreshTokenBase.SignedString(signingKey)

	if err != nil {
		slog.Error("error creating refresh token", "error", err.Error())
		return "", err
	}

	return refreshToken, nil
}

func GetUserIdentityFromAccessToken(tokenString string) *models.User {
	signingKey := []byte(os.Getenv("SECRET_SIGNATURE"))

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return signingKey, nil
	})

	if err != nil || !token.Valid {
		slog.Info("error parsing access token", "error", err.Error())
		return nil
	}

	claims := token.Claims.(models.JWTTokenClaims)

	return &models.User{
		Email: claims.Email,
	}
}

func GetUserIdentityFromRefreshToken(tokenString string) *models.User {
	signingKey := []byte(os.Getenv("SECRET_SIGNATURE"))

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return signingKey, nil
	})

	if err != nil || !token.Valid {
		slog.Info("error parsing refresh token", "error", err.Error())
		return nil
	}

	claims := token.Claims.(models.JWTTokenClaims)

	return &models.User{
		Email: claims.Email,
	}
}

func GenerateNewAccessTokenFromRefreshToken(tokenString string) (string, error) {
	signingKey := []byte(os.Getenv("SECRET_SIGNATURE"))

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return signingKey, nil
	})

	if err != nil || !token.Valid {
		slog.Info("error parsing refresh token", "error", err.Error())
		return "", err
	}

	claims := token.Claims.(models.JWTTokenClaims)

	dbUserRefreshTokenVersion, err := db.GetUserRefreshTokenVersion(claims.Email)

	if err != nil {
		return "", err
	}

	if dbUserRefreshTokenVersion != claims.RefreshTokenVersion {
		slog.Info("refresh token version doesn't match")
		return "", err
	}

	accessToken, err := GenerateAccessToken(claims.Email)

	if err != nil {
		return "", err
	}

	return accessToken, nil
}
