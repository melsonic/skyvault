package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/melsonic/skyvault/blobserver/minio"
	"github.com/melsonic/skyvault/blobserver/types"
)

func chunkSaveHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request body"))
		return
	}
	var data types.BlobData
	err = json.Unmarshal(body, &data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request body format"))
		return
	}
	err = minio.UploadChunk(data.Hash, data.Data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Blob chunk uploaded succesfully!"))
	return
}

func chunkGetHandler(w http.ResponseWriter, r *http.Request) {
	requestQueryMap := r.URL.Query()
	hash := requestQueryMap.Get("hash")
	if hash == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Query parameter 'hash' missing!"))
		return
	}
	content, err := minio.GetChunk(hash)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	blob := types.BlobDataResponse{
		Data:    content,
		Message: "Blob chunk fetched succesfully!",
	}

	response, err := json.Marshal(blob)

	w.WriteHeader(http.StatusOK)
	w.Write(response)
	return
}

func main() {
	minio.InitMinio()
	mux := http.NewServeMux()
	mux.HandleFunc("/chunk/save", chunkSaveHandler)
	mux.HandleFunc("/chunk/get", chunkGetHandler)
	server := &http.Server{
		Addr:           ":8002",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(server.ListenAndServe())
}
