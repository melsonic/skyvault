package db

import (
	"errors"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/melsonic/skyvault/auth/models"
	"github.com/melsonic/skyvault/auth/util"
)

var (
	DBConnPool *pgx.ConnPool
)

func connectDB() error {
	portString := os.Getenv("DB_PORT")
	port, err := strconv.ParseUint(portString, 10, 16)
	if err != nil {
		return err
	}
	dbConfig := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     uint16(port),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Database: os.Getenv("DB_DATABASE"),
		},
		AfterConnect: func(c *pgx.Conn) error {
			slog.Info("DB connected!")
			return nil
		},
	}
	DBConnPool, err = pgx.NewConnPool(dbConfig)
	if err != nil {
		return err
	}
	return nil
}

func setupDB() error {
	_, err := DBConnPool.Exec(`
		CREATE TYPE gender AS ENUM ('male', 'female', 'other');
	`)
	if err != nil {
		slog.Error("error creating gender enum", "error", err.Error())
		return errors.New("error creating gender enum")
	}

	_, err = DBConnPool.Exec(`--sql
		CREATE TABLE IF NOT EXISTS users (
			id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			name TEXT,
			user_gender gender,
			date_of_birth DATE,
			date_created DATE DEFAULT CURRENT_DATE,
			refresh_token_version BIGINT DEFAULT 1
		)
	`)
	if err != nil {
		slog.Error("error creating users table", "error", err.Error())
		return errors.New("error creating users table")
	}

	_, err = DBConnPool.Exec(`--sql
		CREATE TABLE IF NOT EXISTS password_reset_requests (
			id TEXT NOT NULL,
			request_expiry TIMESTAMPTZ NOT NULL,
			email VARCHAR(255) REFERENCES users(email),
			used_flag BOOLEAN DEFAULT FALSE
		)
	`)
	if err != nil {
		slog.Error("error creating password_reset_requests table", "error", err.Error())
		return errors.New("error creating password_reset_requests table")
	}

	return nil
}

func InitDB() error {
	err := connectDB()
	if err != nil {
		return err
	}
	err = setupDB()
	return err
}

func CreateUser(user *models.User) error {

	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return err
	}

	err = DBConnPool.QueryRow(`--sql
		INSERT INTO users (email, password, name, user_gender, date_of_birth, date_created) VALUES ($1, $2, $3, $4, $5, current_timestamp) RETURNING refresh_token_version
	`, user.Email, hashedPassword, user.Name, user.Gender).Scan(&user.RefreshTokenVersion)

	if err != nil {
		slog.Error("error signing up user", "error", err.Error())
		return err
	}

	return nil
}

func GetUserProfile(user *models.User) error {

	err := DBConnPool.QueryRow(`--sql
		SELECT 
			name, user_gender, date_of_birth, date_created
		FROM users 
		WHERE
			email = $1
	`, user.Email).Scan(user.Password, user.Name, user.Gender, user.DateOfBirth, user.DateCreated)

	if err != nil {
		slog.Error("error fetching user details", "error", err.Error())
		return err
	}

	return nil
}

func GetUserName(email string) (string, error) {
	var name string

	err := DBConnPool.QueryRow(`--sql
		SELECT
			name
		FROM users
		WHERE
			email = $1
	`, email).Scan(&name)

	if err != nil {
		slog.Error("error fetching name from users", "error", err.Error())
		return "", err
	}

	return name, nil
}

func GetUserRefreshTokenVersion(email string) (int, error) {
	var refreshTokenVersion int

	err := DBConnPool.QueryRow(`--sql
		SELECT refresh_token_version FROM users WHERE email = $1
	`, email).Scan(&refreshTokenVersion)

	if err != nil {
		slog.Error("error fetching refresh token version", "error", err.Error())
		return -1, err
	}

	return refreshTokenVersion, nil
}

func UpdateRefreshTokenVersion(email string) error {
	var refreshTokenVersion int

	err := DBConnPool.QueryRow(`--sql
		SELECT refresh_token_version FROM users WHERE email = $1
	`, email).Scan(&refreshTokenVersion)

	if err != nil {
		slog.Error("error fetching refresh token version", "error", err.Error())
		return err
	}

	_, err = DBConnPool.Exec(`--sql
		UPDATE users
		SET refresh_token_version = $1
		WHERE email = $2
	`, refreshTokenVersion+1, email)

	if err != nil {
		slog.Error("error updating refresh token version", "error", err.Error())
		return err
	}

	return nil
}

func VerifyUserPasswordWithDBPassword(email string, password string) bool {
	var dbPassword string

	// in case of no user found, the Scan will return ErrNoRows
	err := DBConnPool.QueryRow(`--sql
		SELECT password FROM users WHERE email = $1
	`, email).Scan(&dbPassword)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Error("user doesn't exist", "error", err.Error())
		} else {
			slog.Error("error fetching password from `users`", "error", err.Error())
		}
		return false
	}

	validPassword := util.IsValidPassword(password, dbPassword)

	return validPassword
}

func DeleteUserProfile(email string) error {
	_, err := DBConnPool.Exec(`--sql
		DELETE FROM users
		WHERE email = $1
	`, email)

	if err != nil {
		slog.Error("error deleting user profile.", "error", err.Error())
		return err
	}

	return nil
}

func UpdateUserProfile(realUserIdentity models.User, toUpdateUser models.User) error {
	_, err := DBConnPool.Exec(`--sql
		UPDATE users
		SET
			name = $1, user_gender = $2, date_of_birth = $3
		WHERE
			email = $4
	`, toUpdateUser.Name, toUpdateUser.Gender, toUpdateUser.DateOfBirth, realUserIdentity.Email)

	if err != nil {
		slog.Error("error updating user profile", "error", err.Error())
		return err
	}

	return nil
}

func UpdateUserPassword(email string, password string) error {
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return err
	}

	_, err = DBConnPool.Exec(`--sql
		UPDATE users
		SET
			password = $1
		WHERE
			email = $2
	`, hashedPassword, email)

	if err != nil {
		slog.Error("error updating password", "error", err.Error())
		return err
	}

	return nil
}

func InsertNewPasswordResetRequest(email string) (string, error) {
	var userID int

	err := DBConnPool.QueryRow(`--sql
		SELECT id FROM users WHERE email = $1
	`, email).Scan(userID)

	if err != nil {
		slog.Error("error fetching user identity", "error", err.Error())
		return "", err
	}

	id := uuid.New()

	// expire password reset request id after 1 hour
	_, err = DBConnPool.Exec(`--sql
		INSERT INTO password_reset_requests (id, email, request_expiry) VALUES ($1, $2, $3)
	`, id, email, time.Now().Add(time.Hour))

	if err != nil {
		slog.Error("error inserting new password reset request", "error", err.Error())
		return "", err
	}

	return id.String(), nil
}

func GetPasswordResetRequest(id string) (*models.PasswordResetRequest, error) {
	var PasswordResetRequest models.PasswordResetRequest

	err := DBConnPool.QueryRow(`--sql 
	    SELECT id, email, request_expiry, used_flag FROM password_reset_requests WHERE id = $1
	`, id).Scan(PasswordResetRequest.ID, PasswordResetRequest.Email, PasswordResetRequest.RequestExpiry, PasswordResetRequest.UsedFlag)

	if err != nil {
		slog.Error("error fetching password reset request details", "error", err.Error())
		return nil, err
	}

	return &PasswordResetRequest, nil
}

func MarkPasswordResetRequestAsUsed(requestID string) error {
	_, err := DBConnPool.Exec(`--sql
		UPDATE password_reset_requests
		SET
			used_flag = TRUE
		WHERE 
			id = $1
	`, requestID)

	if err != nil {
		slog.Error("error setting used flag", "error", err.Error(), "request_id", requestID)
		return err
	}

	return nil
}
