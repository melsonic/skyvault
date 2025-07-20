package db

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	"github.com/joho/godotenv"
	"github.com/melsonic/skyvault/metadata/types"
	"github.com/melsonic/skyvault/metadata/util"
)

var (
	DBConnPool  *pgx.ConnPool
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
		CREATE TABLE IF NOT EXISTS NODE (
			ID bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			FOLDER boolean NOT NULL,
			NAME text NOT NULL,
			PARENT_FOLDER bigint,
			CREATED_AT timestamptz NOT NULL,
			LAST_ACCESS timestamptz NOT NULL,
			LAST_MODIFIED timestamptz NOT NULL
		)
	`)
	if err != nil {
		slog.Error("error creating NODE table", "error", err.Error())
		return errors.New("error creating node table")
	}

	_, err = DBConnPool.Exec(`
		CREATE TABLE IF NOT EXISTS FILE_METADATA (
			ID bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			FILE_TYPE text NOT NULL,
			FILE_SIZE bigint NOT NULL,
			HASH_IDS text[] NOT NULL,
			NODE_ID bigint references NODE(ID)
		)
	`)
	if err != nil {
		slog.Error("error creating FILE_METADATA table", "error", err.Error())
		return errors.New("error creating FILE_METADATA table")
	}

	_ = DBConnPool.QueryRow("SELECT ID FROM NODE WHERE NAME = $1", ROOT_NAME).Scan(&ROOT_FOLDER)
	if ROOT_FOLDER == -1 {
		err = DBConnPool.QueryRow(`INSERT INTO NODE (FOLDER, NAME, CREATED_AT, LAST_ACCESS, LAST_MODIFIED) VALUES ($1, $2, current_timestamp, current_timestamp, current_timestamp) RETURNING ID`, true, ROOT_NAME).Scan(&ROOT_FOLDER)
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

func SaveMetadata(data *types.Metadata) (int, error) {
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

		if path[index] == "" {
			continue
		}

		err := DBConnPool.QueryRow(`
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
	for ; index < len(path); index++ {

		if path[index] == "" {
			continue
		}

		err := DBConnPool.QueryRow(`
			INSERT INTO NODE (FOLDER, NAME, PARENT_FOLDER, CREATED_AT, LAST_ACCESS, LAST_MODIFIED) 
			VALUES 
				($1, $2, $3, current_timestamp, current_timestamp, current_timestamp) 
			RETURNING ID
		`, true, path[index], folderID).Scan(&currentFolderID)
		if err != nil {
			slog.Error("error inserting node", "error", err.Error(), "nodename", path[index])
			return -1, errors.New("error saving node")
		}

		folderID = currentFolderID
	}
	if data.IsFolder {
		return folderID, nil
	}
	// Below code runs only when the input is a file
	hashes := util.FormatHashedChunks(data.Hashes)

	var nodeID int
	err = DBConnPool.QueryRow(`
		INSERT INTO NODE (
			FOLDER, NAME, PARENT_FOLDER, CREATED_AT, LAST_ACCESS, LAST_MODIFIED
		) 
		VALUES 
			($1, $2, $3, current_timestamp, current_timestamp, current_timestamp)
		RETURNING ID, CREATED_AT, LAST_ACCESS, LAST_MODIFIED
	`, data.IsFolder, data.FileName, folderID).Scan(&nodeID, &data.CreatedAt, &data.LastAccess, &data.LastModified)
	if err != nil {
		slog.Error("error saving file", "error", err.Error(), "filename", data.FileName)
		return -1, errors.New("error saving file")
	}

	_, err = DBConnPool.Exec(`
		INSERT INTO FILE_METADATA (
			FILE_TYPE, FILE_SIZE, NODE_ID, HASH_IDS
		) 
		VALUES 
		(
			$1, $2, $3, $4
		)
	`, fileExtension, data.FileSize, nodeID, hashes)
	if err != nil {
		slog.Error("error inserting node", "error", err.Error())
		return -1, errors.New("error saving node")
	}
	return nodeID, nil
}

func FetchMetadata(nodeID string) (*types.Metadata, error) {
	var hashIDs pgtype.TextArray
	var fileNodeData types.Metadata
	err := DBConnPool.QueryRow(`
		SELECT 
			NODE.CREATED_AT, NODE.LAST_ACCESS, NODE.LAST_MODIFIED, FILE_METADATA.FILE_SIZE, FILE_METADATA.HASH_IDS, NODE.NAME 
		FROM 
			NODE, FILE_METADATA 
		WHERE 
			NODE.ID = FILE_METADATA.NODE_ID AND NODE.ID=$1
	`, nodeID).Scan(&fileNodeData.CreatedAt, &fileNodeData.LastAccess, &fileNodeData.LastModified, &fileNodeData.FileSize, &hashIDs, &fileNodeData.FileName)
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

// In case of delete, we need to delete all the child folders/files
func DeleteMetadata(nodeID string) error {
	id, err := strconv.ParseInt(nodeID, 10, 64)
	if err != nil {
		slog.Error("invalid nodeID", "id", nodeID, "error", err.Error())
		return errors.New("invalid node id")
	}
	queue := []int64{id}
	nodeIDs := []int64{id}

	// bfs
	for len(queue) > 0 {
		currentNodeID := queue[0]
		queue = queue[1:]
		rows, err := DBConnPool.Query(`
			SELECT 
				ID
			FROM 
				NODE
			WHERE
				PARENT_FOLDER = $1
		`, currentNodeID)
		if err != nil {
			slog.Error("error travelling all nodes", "error", err.Error())
			return errors.New("error deleting node")
		}
		defer rows.Close()

		var id int64
		for rows.Next() {
			if err := rows.Scan(&id); err != nil {
				slog.Error("error fetching id from rows", "error", err.Error())
				return errors.New("error deleting node")
			}
			queue = append(queue, id)
			nodeIDs = append(nodeIDs, id)
		}

		if err = rows.Err(); err != nil {
			slog.Error("error fetching id from rows end", "error", err.Error())
			return errors.New("error deleting node")
		}
	}

	for i := range nodeIDs {
		_, err = DBConnPool.Exec(`
			DELETE FROM
				FILE_METADATA
			WHERE
				NODE_ID = $1
		`, nodeIDs[i])

		if err != nil {
			slog.Error("error deleting row in file_metadata", "error", err.Error(), "node_id", nodeIDs[i])
			return errors.New("error deleting file_metadata")
		}

		_, err = DBConnPool.Exec(`
			DELETE FROM
				NODE
			WHERE
				ID = $1
		`, nodeIDs[i])

		if err != nil {
			slog.Error("error deleting row node", "error", err.Error(), "id", nodeIDs[i])
			return errors.New("error deleting node")
		}

	}

	return nil
}

// In case while deleting a folder, there might be some error
// and some orphan nodes might be present in node table
// this method runs periodically to remove those orphan nodes
func CleanOrphanNodes(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)

	go func() {
		for {
			select {
			case <-ticker.C:
				// clean up the table
				// TODO: optimize the query/approach
				_, err := DBConnPool.Exec(`DELETE FROM NODE WHERE PARENT_FOLDER NOT IN (SELECT ID FROM NODE)`)
				if err != nil {
					slog.Error("error fetching node ids in gorouting", "error", err.Error())
					continue
				}

			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()

}
