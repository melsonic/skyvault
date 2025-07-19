package main

import (
	"fmt"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/melsonic/skyvault/metadata/db"
	"github.com/melsonic/skyvault/metadata/types"
)

func metadataSaveHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error reading request body"))
		return
	}
	var data types.Metadata
	err = json.Unmarshal(body, &data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request body"))
		return
	}
	// Perform Operation to save Metadata
	nodeID, err := db.SaveMetadata(data)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	data.FileNodeId = fmt.Sprint(nodeID)
	jsonData, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("marshal error"))
		return
	}
	// write json data to writer
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func metadataFetchHandler(w http.ResponseWriter, r *http.Request) {
	nodeid := r.PathValue("nodeid")
	// Perform operation to fetch hashes from Database
	nodeMetaData, err := db.FetchMetadata(nodeid)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error fetching file metadata"))
		return
	}
	response, err := json.Marshal(nodeMetaData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func main() {
	err := db.InitDB()
	if err != nil {
		slog.Error("error in db setup", "error", err.Error())
		return
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /metadata/{nodeid}", metadataFetchHandler)
	mux.HandleFunc("POST /metadatas", metadataSaveHandler)
	server := &http.Server{
		Addr:           ":8001",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(server.ListenAndServe())
}
