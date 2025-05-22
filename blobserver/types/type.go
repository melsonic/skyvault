package types

type BlobData struct {
	Hash string
	Data []byte
}

type BlobDataResponse struct {
	Data    []byte
	Message string
}
