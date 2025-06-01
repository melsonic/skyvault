package types

type BlobData struct {
	Hash string `json:"hash"`
	Data []byte `json:"data"`
}

type BlobDataResponse struct {
	Data    []byte `json:"data"`
	Message string `json:"message"`
}
