package main

import (
	"log"
	"net/http"
	"summersea.top/filetransfer"
)

func main() {
	server := filetransfer.NewFileServer()

	err := http.ListenAndServe(":80", server)
	if err != nil {
		log.Fatalf("could not listen on port 80: %v", err)
	}
}
