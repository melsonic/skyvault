package types

import "time"

// FilePath will not contain the filename
// If IsFolder is true, FileName & Hashes will be empty
type Metadata struct {
	FileNodeId   string    `json:"nodeid"`
	FileName     string    `json:"filename"`
	FilePath     string    `json:"filepath"`
	IsFolder     bool      `json:"is_folder"`
	Hashes       []string  `json:"hashes"`
	FileSize     int       `json:"filesize"`
	CreatedAt    time.Time `json:"created_at"`
	LastAccess   time.Time `json:"last_access"`
	LastModified time.Time `json:"last_modified"`
}
