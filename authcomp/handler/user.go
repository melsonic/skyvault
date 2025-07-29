package handler

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/melsonic/skyvault/auth/db"
	"github.com/melsonic/skyvault/auth/models"
	"github.com/melsonic/skyvault/auth/util"
)

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := ctx.Value("user").(models.User)

	err := db.GetUserProfile(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error fetching user profile."))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("user profile fetched."))
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	realUserIdentity := ctx.Value("user").(models.User)

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

	if data.Name == "" {
		slog.Error("empty required field name")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("empty required field name"))
		return
	}

	if !util.IsValidEmail(data.Email) {
		slog.Error("invalid email address in profile update", "email", data.Email)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid email address"))
		return
	}

	if data.Gender != "" && !util.IsValidGender(data.Gender) {
		slog.Error("invalid gender in profile update", "gender", data.Gender)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid gender"))
		return
	}

	err = db.UpdateUserProfile(realUserIdentity, data)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error updating user profile."))
		return
	}
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("user profile updated successfully."))
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := ctx.Value("user").(models.User)

	err := db.DeleteUserProfile(user.Email)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error deleting user profile."))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("user profile deleted successfully."))
}
