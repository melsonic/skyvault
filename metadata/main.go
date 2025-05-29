package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	err = db.SaveMetadata(data)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File metadata saved successfully!"))
	return
}

func metadataGetHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request body"))
		return
	}
	var data types.Metadata
	err = json.Unmarshal(body, &data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request body"))
		return
	}
	// Perform operation to fetch hashes from Database
	hashes, err := db.FetchMetadata(data.FileName)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error fetching file metadata"))
		return
	}
	md := types.Metadata{
		Hashes: hashes,
	}
	response, err := json.Marshal(md)
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func main() {
	err := db.InitDB()
	if err != nil {
		fmt.Println("InitDB Error : ", err.Error())
		return
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/metadata/save", metadataSaveHandler)
	mux.HandleFunc("/metadata/get", metadataGetHandler)
	server := &http.Server{
		Addr:           ":8001",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(server.ListenAndServe())
}
