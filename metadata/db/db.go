package db

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
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
	slog.Info("database connected!")
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
		slog.Error("error creating NODE_MD table", "error", err.Error())
		return errors.New("error creating node_md table")
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
		slog.Error("error creating NODE table", "error", err.Error())
		return errors.New("error creating node table")
	}
	_ = DBConn.QueryRow("SELECT ID FROM NODE WHERE NAME = $1", ROOT_NAME).Scan(&ROOT_FOLDER)
	if ROOT_FOLDER == -1 {
		err = DBConn.QueryRow(`INSERT INTO NODE (FOLDER, NAME) VALUES ($1, $2) RETURNING ID`, true, ROOT_NAME).Scan(&ROOT_FOLDER)
		if err != nil {
			slog.Error("error inserting root node", "error", err.Error())
			return errors.New("error inserting root node")
		}
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

func SaveMetadata(data types.Metadata) (int, error) {
	fileExtension, err := util.GetFileExtension(data.FileName)
	if !data.IsFolder && err != nil {
		slog.Error("unsupported file format")
		return -1, err
	}
	path := strings.Split(data.FilePath, "/")
	folderID := ROOT_FOLDER
	var index int = 0
	var currentFolderID int
	for index = 0; index < len(path); index++ {
		err := DBConn.QueryRow(`
			SELECT 
				ID 
			FROM 
				NODE 
			WHERE 
				NAME = $1
		`, path[index]).Scan(&currentFolderID)
		if err != nil {
			break
		}
		folderID = currentFolderID
	}
	for index < len(path) {
		err := DBConn.QueryRow(`
			INSERT INTO NODE (FOLDER, NAME, PARENT_FOLDER) 
			VALUES 
				($1, $2, $3) RETURNING ID
		`, true, path[index], folderID).Scan(&currentFolderID)
		if err != nil {
			slog.Error("error inserting node", "error", err.Error(), "nodename", path[index])
			return -1, errors.New("error saving node")
		}
		folderID = currentFolderID
		index++
	}
	if data.IsFolder {
		return folderID, nil
	}
	// Below code runs only when the input is a file
	var nodeMetaDataID int = -1
	err = DBConn.QueryRow(`
		INSERT INTO NODE_MD (
			FILE_TYPE, CREATED_AT, LAST_ACCESS, 
			LAST_MODIFIED, FILE_SIZE
		) 
		VALUES 
		(
			$1, current_timestamp, current_timestamp, 
			current_timestamp, $2
		) RETURNING ID
	`, fileExtension, data.FileSize).Scan(&nodeMetaDataID)
	if err != nil {
		slog.Error("error inserting node", "error", err.Error())
		return -1, errors.New("error saving node")
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
	var nodeID int
	err = DBConn.QueryRow(`
		INSERT INTO NODE (
			FOLDER, NAME, PARENT_FOLDER, FILE_MD, 
			HASH_IDS
		) 
		VALUES 
			($1, $2, $3, $4, $5)
		RETURNING ID
	`, data.IsFolder, data.FileName, folderID, nodeMetaDataID, hashes).Scan(&nodeID)
	if err != nil {
		slog.Error("error saving file", "error", err.Error(), "filename", data.FileName)
		return -1, errors.New("error saving file")
	}
	return nodeID, nil
}

func FetchMetadata(fileNodeId string) (*types.Metadata, error) {
	var hashIDs pgtype.TextArray
	var fileNodeData types.Metadata
	err := DBConn.QueryRow(`
		SELECT 
			NODE_MD.CREATED_AT, NODE_MD.LAST_ACCESS, NODE_MD.LAST_MODIFIED, NODE_MD.FILE_SIZE, NODE.HASH_IDS, NODE.NAME 
		FROM 
			NODE, NODE_MD 
		WHERE 
			NODE.ID=$1
	`, fileNodeId).Scan(&fileNodeData.CreatedAt, &fileNodeData.LastAccess, &fileNodeData.LastModified, &fileNodeData.FileSize, &hashIDs, &fileNodeData.FileName)
	if err != nil {
		slog.Error("error in querying node db", "error", err.Error())
		return nil, errors.New("error fetching data from db")
	}
	var hashes []string
	for i := range hashIDs.Elements {
		hashes = append(hashes, hashIDs.Elements[i].String)
	}
	fileNodeData.Hashes = hashes
	return &fileNodeData, nil
}
