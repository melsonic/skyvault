package types

// FilePath will not contain the filename
// If IsFolder is true, FileName & Hashes will be empty
type Metadata struct {
	FileName string   `json:"filename"`
	FilePath string   `json:"filepath"`
	IsFolder bool     `json:"is_folder"`
	Hashes   []string `json:"hashes"`
}
