package db

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	"github.com/joho/godotenv"
	"github.com/melsonic/skyvault/metadata/types"
)

var (
	DBConn      *pgx.Conn
	ROOT_FOLDER = -1
	ROOT_NAME   = "root"
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
		CREATE TABLE IF NOT EXISTS NODE (
			ID bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			FOLDER boolean NOT NULL,
			NAME text NOT NULL,
			PARENT_FOLDER bigint,
			HASH_IDS text[]
		)
	`)
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
	path := strings.Split(data.FilePath, "/")
	var parentFolderId int
	var index int = 0
	for index = 0; index < len(path); index++ {
		var tempParentId int
		err := DBConn.QueryRow("SELECT ID FROM NODE WHERE NAME = $1", path[index]).Scan(&tempParentId)
		if err != nil {
			break
		}
		parentFolderId = tempParentId
	}
	for index < len(path) {
		var tempParentId int
		err := DBConn.QueryRow("INSERT INTO NODE (FOLDER, NAME, PARENT_FOLDER) VALUES ($1, $2, $2) RETURNING ID", true, path[index], parentFolderId).Scan(&tempParentId)
		if err != nil {
			return err
		}
		parentFolderId = tempParentId
		index++
	}
	if data.IsFolder {
		return nil
	}
	hashes := "{"
	for i := range data.Hashes {
		hashes += fmt.Sprintf(`"%s"`, data.Hashes[i])
		if i == len(data.Hashes)-1 {
			break
		}
		hashes += ","
	}
	hashes += "}"
	_, err := DBConn.Exec(`INSERT INTO NODE (FOLDER, NAME, PARENT_FOLDER, HASH_IDS) VALUES ($1, '$2', $3, $4)`, data.IsFolder, data.FileName, parentFolderId, hashes)
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
