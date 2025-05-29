package db

import (
	"fmt"
	"log"
	"os"
	"strconv"
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
	portString := os.Getenv("DB_PORT")
	port, err := strconv.ParseUint(portString, 10, 16)
	if err != nil {
		return err
	}
	dbConfig := pgx.ConnConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     uint16(port),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
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
	if err != nil {
		return err
	}
	err = DBConn.QueryRow("SELECT ID FROM NODE WHERE NAME = $1", ROOT_NAME).Scan(&ROOT_NAME)
	if err != nil {
		err = DBConn.QueryRow(`INSERT INTO NODE (FOLDER, NAME) VALUES ($1, $2) RETURNING ID`, true, ROOT_NAME).Scan(&ROOT_FOLDER)
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

func SaveMetadata(data types.Metadata) error {
	path := strings.Split(data.FilePath, "/")
	parentFolderId := ROOT_FOLDER
	var index int = 0
	for index = 0; index < len(path); index++ {
		var tempParentId int
		err := DBConn.QueryRow("SELECT ID FROM NODE WHERE NAME = $1", path[index]).Scan(&tempParentId)
		if err != nil {
			break
		}
		parentFolderId = tempParentId
	}
	fmt.Println(index)
	for index < len(path) {
		var tempParentId int
		fmt.Println(path[index])
		err := DBConn.QueryRow("INSERT INTO NODE (FOLDER, NAME, PARENT_FOLDER) VALUES ($1, $2, $3) RETURNING ID", true, path[index], parentFolderId).Scan(&tempParentId)
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
	_, err := DBConn.Exec(`INSERT INTO NODE (FOLDER, NAME, PARENT_FOLDER, HASH_IDS) VALUES ($1, $2, $3, $4)`, data.IsFolder, data.FileName, parentFolderId, hashes)
	return err
}

func FetchMetadata(fileName string) ([]string, error) {
	var hashIDs pgtype.TextArray
	err := DBConn.QueryRow("SELECT HASH_IDS FROM NODE WHERE NAME=$1", fileName).Scan(&hashIDs)
	if err != nil {
		return nil, err
	}
	var hashes []string
	for i := range hashIDs.Elements {
		hashes = append(hashes, hashIDs.Elements[i].String)
	}
	return hashes, nil
}
