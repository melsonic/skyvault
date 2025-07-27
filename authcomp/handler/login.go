package handler

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/melsonic/skyvault/auth/db"
	jwtauth "github.com/melsonic/skyvault/auth/jwt"
	"github.com/melsonic/skyvault/auth/models"
	"github.com/melsonic/skyvault/auth/util"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("error reading request body", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to read request body"))
		return
	}

	var data models.User
	err = json.Unmarshal(requestBody, &data)
	if err != nil {
		slog.Error("error in json unmarshal", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error unmarshalling json value"))
		return
	}

	if !util.IsValidEmail(data.Email) {
		slog.Error("invalid email address", "email", data.Email)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid email address"))
		return
	}

	userPresent := db.VerifyUserProfile(data.Email, data.Password)

	if !userPresent {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user not present"))
		return
	}

	data.Name, err = db.GetUserName(data.Email);

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error fetching user details"))
		return
	}

	// access token & refresh token generation
	refreshToken, err := jwtauth.GenerateRefreshToken(data.Email, data.Name, data.RefreshTokenVersion)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error generating auth token[1]"))
		return
	}

	accessToken, err := jwtauth.GenerateAccessToken(data.Email, data.Name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error generating auth token[2]"))
		return
	}

	authResponse := models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	authResponseJsonData, err := json.Marshal(authResponse)

	if err != nil {
		slog.Error("error in json marshal", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("please try again later!"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(authResponseJsonData)
}
