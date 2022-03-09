package main

import (
	"log"
	"summersea.top/filetransfer"
)

func main() {
	store := filetransfer.CreateStoreByConfig()
	adapter := filetransfer.NewFileTranDataAdapter(store)
	server := filetransfer.NewFileServer(adapter)

	err := server.Run(":8080")
	if err != nil {
		log.Fatalf("could not listen on port 80: %v", err)
	}
}
