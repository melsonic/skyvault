package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/melsonic/skyvault/auth/models"
)

func AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := ctx.Value("user").(*models.User)

	userJsonData, err := json.Marshal(user)

	if err != nil {
		slog.Error("error in json marshal", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("please try again later!"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userJsonData)
}
