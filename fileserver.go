package filetransfer

import (
	"encoding/json"
	"fmt"
	"github.com/kirinlabs/utils/str"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

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

func (fs *FileServer) uploadInitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	fs.handleUploadInit(w, r)
}

func (fs *FileServer) handleUploadInit(w http.ResponseWriter, r *http.Request) {
	uploadInitBody, err := fs.extractBody(r)
	if err != nil {
		log.Printf("problem extract request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !isUploadInitReqBodyValid(*uploadInitBody) {
		log.Printf("got unsupported request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uid := uuid.NewV4()
	_, err = w.Write([]byte(uid.String()))
	if err != nil {
		log.Printf("problem write data to response")
	}
}

func (fs *FileServer) extractBody(r *http.Request) (*UploadInitReqBody, error) {
	var uploadInitBody UploadInitReqBody
	if r.Body == nil {
		err := fmt.Errorf("got nil request body")
		return nil, err
	}
	err := json.NewDecoder(r.Body).Decode(&uploadInitBody)
	if err != nil {
		err = fmt.Errorf("problem decode request body: %v", err)
		return nil, err
	}
	return &uploadInitBody, nil
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
