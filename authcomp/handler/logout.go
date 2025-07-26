package handler

import (
	"net/http"

	"github.com/melsonic/skyvault/auth/db"
	"github.com/melsonic/skyvault/auth/models"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := ctx.Value("user").(models.User)

	err := db.UpdateRefreshTokenVersion(user.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error logging out user"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("user logged out successfully!"))
}
