package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/melsonic/skyvault/auth/db"
	"github.com/melsonic/skyvault/auth/models"
	"github.com/melsonic/skyvault/auth/util"
	"github.com/wneessen/go-mail"
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

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("error reading request body", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error reading request body"))
		return
	}

	var passwordResetRequest models.PasswordResetRequest
	err = json.Unmarshal(requestBody, &passwordResetRequest)
	if err != nil {
		slog.Error("error in json unmarshal", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error reading request body"))
		return
	}

	if !util.IsValidEmail(passwordResetRequest.Email) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid email"))
		return
	}

	passwordResetRequest.ID, err = db.InsertNewPasswordResetRequest(passwordResetRequest.Email)
	if err != nil {
		slog.Error("error generating uuid for password reset request", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error posting password reset request"))
		return
	}

	// Using MailHog for local SMTP server
	email := mail.NewMsg()

	err = email.From("support@skyvault.com")
	if err != nil {
		slog.Error("error setting email sender", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error sending email"))
		return
	}

	err = email.To(passwordResetRequest.Email)
	if err != nil {
		slog.Error("error setting email sender", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error sending email"))
		return
	}

	email.Subject("Password Reset - SkyVault")

	emailTemplate, err := template.ParseFiles("../template/email.tmpl")
	if err != nil {
		slog.Error("error parsing tmpl file", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error sending email"))
		return
	}
	data := models.EmailTemplateData{
		ResetURL: fmt.Sprintf("http://localhost:8003/user/updatepassword/%s", passwordResetRequest.ID),
	}

	err = email.SetBodyHTMLTemplate(emailTemplate, &data)
	if err != nil {
		slog.Error("error setting email body", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error sending email"))
		return
	}

	emailClient, err := mail.NewClient("0.0.0.0", mail.WithPort(1025))
	if err != nil {
		slog.Error("error creating email client", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error sending email"))
		return
	}
	emailClient.SetTLSPolicy(mail.TLSOpportunistic)

	if err = emailClient.DialAndSend(email); err != nil {
		slog.Error("error sending email", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error sending email"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("email sent"))
}

func UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	passwordResetID := r.PathValue("hashid")

	passwordResetRequest, err := db.GetPasswordResetRequest(passwordResetID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid password reset request"))
		return
	}

	if passwordResetRequest.RequestExpiry.Before(time.Now()) {
		slog.Error("password reset request expired")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid password reset request"))
		return
	}

	if passwordResetRequest.UsedFlag {
		slog.Error("password reset request id is already used")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid password reset request"))
		return
	}

	// Now read the new password from request body
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("error reading request body", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid password reset request"))
		return
	}

	var passwordResetUpdateRequest models.PasswordResetUpdateRequest
	err = json.Unmarshal(requestBody, &passwordResetUpdateRequest)
	if err != nil {
		slog.Error("error in json unmarshal", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid password reset request"))
		return
	}

	// Compare old password with current db password
	ok := db.VerifyUserPasswordWithDBPassword(passwordResetRequest.Email, passwordResetUpdateRequest.OldPassword)
	if !ok {
		slog.Error("old password doesn't match")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid password reset request"))
		return
	}

	// Update user password with new password
	err = db.UpdateUserPassword(passwordResetRequest.Email, passwordResetUpdateRequest.NewPassword)
	if err != nil {
		slog.Error("error updating db password", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid password reset request"))
		return
	}

	// Mark password reset request as used
	err = db.MarkPasswordResetRequestAsUsed(passwordResetRequest.ID)
	if err != nil {
		slog.Error("error marking password reset request as used", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid password reset request"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("password updated successfully"))
}
