package filetransfer

import (
	"encoding/json"
	"github.com/kirinlabs/utils/str"
	"log"
	"net/http"
)

type Resource struct {
	Address string  `json:"address"`
	Port    int     `json:"port"`
	Account Account `json:"account"`
}

type Account struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UploadInitReqBody struct {
	Resource Resource `json:"resource"`
	Path     string   `json:"path"`
}

type FileServer struct {
	http.Handler
}

func NewFileServer() *FileServer {
	fileServer := &FileServer{}
	router := http.NewServeMux()
	router.Handle("/file/upload/initialization", http.HandlerFunc(fileServer.uploadInitHandler))

	fileServer.Handler = router
	return fileServer
}

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func (fs *FileServer) uploadInitHandler(w http.ResponseWriter, r *http.Request) {
	var uploadInitBody UploadInitReqBody
	err := json.NewDecoder(r.Body).Decode(&uploadInitBody)
	if err != nil {
		log.Printf("problem decode request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !isUploadInitReqBodyValid(uploadInitBody) {
		log.Printf("got unsupported request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func isUploadInitReqBodyValid(body UploadInitReqBody) bool {
	if !str.StartsWith(body.Path, "/") {
		return false
	}
	resource := body.Resource
	if resource.Port <= 0 || resource.Port > 65535 {
		return false
	}
	if resource.Address == "" {
		return false
	}
	account := resource.Account
	if account.Name == "" {
		return false
	}
	if account.Password == "" {
		return false
	}
	return true
}
