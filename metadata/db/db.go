package db

import (
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	"github.com/joho/godotenv"
	"github.com/melsonic/skyvault/metadata/types"
)

var (
	DBConn *pgx.Conn
)

func connectDB() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbConfig := pgx.ConnConfig{
		Host:     os.Getenv("HOST"),
		Port:     5432,
		User:     os.Getenv("USER"),
		Password: os.Getenv("PASSWORD"),
	}
	DBConn, err = pgx.Connect(dbConfig)
	if err != nil {
		return err
	}
	return nil
}

func setupDB() error {
	_, err := DBConn.Exec(`
		CREATE TABLE IF NOT EXISTS FILE (
			ID bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			FILE_NAME text NOT NULL,
			HASH_IDS text[]
		)
	`)
	if err != nil {
		return err
	}
	return err
}

func InitDB() error {
	err := connectDB()
	if err != nil {
		return err
	}
	err = setupDB()
	return err
}

func SaveMetadata(data *types.Metadata) error {
	hashes := "{"
	for i := range data.Hashes {
		hashes += fmt.Sprintf(`"%s"`, data.Hashes[i])
		if i == len(data.Hashes)-1 {
			break
		}
		hashes += ","
	}
	hashes += "}"
	_, err := DBConn.Exec(`INSERT INTO FILE VALUES ('$1', '$2')`, data.FileName, hashes)
	return err
}

func FetchMetadata(fileName string) ([]string, error) {
	var hashIDs pgtype.TextArray
	err := DBConn.QueryRow("SELECT HASH_IDS FROM FILE WHERE FILE_NAME=$1", fileName).Scan(&hashIDs)
	if err != nil {
		return nil, err
	}
	var hashes []string
	for i := range hashIDs.Elements {
		hashes = append(hashes, hashIDs.Elements[i].String)
	}
	return hashes, nil
}
