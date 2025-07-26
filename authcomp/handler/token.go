package handler

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	jwtauth "github.com/melsonic/skyvault/auth/jwt"
	"github.com/melsonic/skyvault/auth/models"
)

func NewTokenHandler(w http.ResponseWriter, r *http.Request) {
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("error parsing form", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error parsing input form"))
		return
	}

	var requestData models.NewAccessTokenRequest
	err = json.Unmarshal(requestBody, &requestData)
	if err != nil {
		slog.Error("error in json unmarshal", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error unmarshalling json value"))
		return
	}

	accessToken, err := jwtauth.GenerateNewAccessTokenFromRefreshToken(requestData.RefreshToken)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error generating access token"))
		return
	}

	authResponse := models.NewAccessTokenResponse{
		AccessToken: accessToken,
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
