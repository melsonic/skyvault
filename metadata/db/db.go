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
	"github.com/melsonic/skyvault/metadata/util"
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
		Database: os.Getenv("DB_DATABASE"),
	}
	DBConn, err = pgx.Connect(dbConfig)
	if err != nil {
		return err
	}
	return nil
}

func setupDB() error {
	_, err := DBConn.Exec(`
		CREATE TABLE IF NOT EXISTS NODE_MD (
			ID bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			FILE_TYPE text NOT NULL,
			CREATED_AT timestamptz NOT NULL,
			LAST_ACCESS timestamptz NOT NULL,
			LAST_MODIFIED timestamptz NOT NULL,
			FILE_SIZE bigint NOT NULL
		)
	`)
	if err != nil {
		return err
	}
	_, err = DBConn.Exec(`
		CREATE TABLE IF NOT EXISTS NODE (
			ID bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			FOLDER boolean NOT NULL,
			NAME text NOT NULL,
			PARENT_FOLDER bigint,
			FILE_MD bigint references NODE_MD(ID),
			HASH_IDS text[]
		)
	`)
	if err != nil {
		return err
	}
	_ = DBConn.QueryRow("SELECT ID FROM NODE WHERE NAME = $1", ROOT_NAME).Scan(&ROOT_FOLDER)
	if ROOT_FOLDER == -1 {
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
	// Below code runs only when the input is a file
	fileExtension, err := util.GetFileExtension(data.FileName)
	if err != nil {
		return err
	}
	var nodeMetaDataId int = -1
	err = DBConn.QueryRow("INSERT INTO NODE_MD (FILE_TYPE, CREATED_AT, LAST_ACCESS, LAST_MODIFIED, FILE_SIZE) VALUES ($1, current_timestamp, current_timestamp, current_timestamp, $2) RETURNING ID", fileExtension, data.FileSize).Scan(&nodeMetaDataId)
	if err != nil {
		return err
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
	_, err = DBConn.Exec(`INSERT INTO NODE (FOLDER, NAME, PARENT_FOLDER, FILE_MD, HASH_IDS) VALUES ($1, $2, $3, $4, $5)`, data.IsFolder, data.FileName, parentFolderId, nodeMetaDataId, hashes)
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
