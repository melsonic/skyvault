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
	hash := r.PathValue("hash")
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
	data.Hash = hash
	if data.Hash == "" || len(data.Data) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing input data!"))
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
}

func chunkGetHandler(w http.ResponseWriter, r *http.Request) {
	hash := r.PathValue("hash")
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
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("marshal error"))
		return
	}

	// write json data to writer
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func chunkDeleteHandler(w http.ResponseWriter, r *http.Request) {
	hash := r.PathValue("hash")
	err := minio.DeleteChunk(hash)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Blob chunk deleted succesfully!"))
}

func main() {
	minio.InitMinio()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /chunk/{hash}", chunkGetHandler)
	mux.HandleFunc("POST /chunk/{hash}", chunkSaveHandler)
	mux.HandleFunc("DELETE /chunk/{hash}", chunkDeleteHandler)
	server := &http.Server{
		Addr:           ":8002",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(server.ListenAndServe())
}
