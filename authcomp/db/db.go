package db

import (
	"errors"
	"log/slog"
	"os"
	"strconv"

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

	_, err = DBConnPool.Exec(`
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
		return errors.New("error creating node table")
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

	err = DBConnPool.QueryRow(`
		INSERT INTO users (email, password, name, user_gender, date_of_birth, date_created) VALUES ($1, $2, $3, $4, $5, current_timestamp) RETURNING refresh_token_version
	`, user.Email, hashedPassword, user.Name, user.Gender).Scan(&user.RefreshTokenVersion)

	if err != nil {
		slog.Error("error signing up user", "error", err.Error())
		return err
	}

	return nil
}

func GetUserData(user *models.User) error {

	err := DBConnPool.QueryRow(`
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

	err := DBConnPool.QueryRow(`
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

	err := DBConnPool.QueryRow(`
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

	err := DBConnPool.QueryRow(`
		SELECT refresh_token_version FROM users WHERE email = $1
	`, email).Scan(&refreshTokenVersion)

	if err != nil {
		slog.Error("error fetching refresh token version", "error", err.Error())
		return err
	}

	_, err = DBConnPool.Exec(`
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

func VerifyUserProfile(email string, password string) bool {
	var dbPassword string

	err := DBConnPool.QueryRow(`
		SELECT password FROM users WHERE email = $1
	`, email).Scan(&dbPassword)

	if err != nil {
		slog.Error("error fetching password from `users`", "error", err.Error())
		return false
	}

	validPassword := util.IsValidPassword(password, dbPassword)

	return validPassword
}
