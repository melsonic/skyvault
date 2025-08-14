package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/melsonic/skyvault/auth/db"
	"github.com/melsonic/skyvault/auth/models"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := ctx.Value("user").(*models.User)

	err := db.UpdateRefreshTokenVersion(user.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error logging out user"))
		return
	}

	responseMessage := models.ResponseMessage{
		Message: "user logged out successfully",
	}

	jsonResponseMessage, err := json.Marshal(responseMessage)

	if err != nil {
		slog.Error("error marshaling json", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error sending email"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponseMessage)
}
