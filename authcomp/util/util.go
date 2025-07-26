package util

import (
	"log/slog"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	BCRYPT_COST  = 10
	EMAIL_REGEXP = `^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`
)

func IsValidGender(gender string) bool {
	lowerCaseGender := strings.ToLower(gender)
	return lowerCaseGender == "male" || lowerCaseGender == "female" || lowerCaseGender == "other"
}

func IsValidEmail(email string) bool {
	rexp := regexp.MustCompile(EMAIL_REGEXP)
	return rexp.Match([]byte(email))
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), BCRYPT_COST)
	if err != nil {
		slog.Error("error hashing password", "password", password)
		return "", err
	}
	return string(hashedPassword), nil
}

func IsValidPassword(userPassword string, dbPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(userPassword))
	if err != nil {
		slog.Error("error matching password", "error", err.Error())
		return false
	}
	return true
}
