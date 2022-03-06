package main

import (
	"log"
	"net/http"
	"summersea.top/filetransfer"
)

func main() {
	store := filetransfer.NewMemoryStore()
	adapter := filetransfer.NewFileTranDataAdapter(store)
	server := filetransfer.NewFileServer(adapter)

	err := http.ListenAndServe(":8080", server)
	if err != nil {
		log.Fatalf("could not listen on port 80: %v", err)
	}
}
