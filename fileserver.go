package filetransfer

import (
	"encoding/json"
	"fmt"
	"github.com/kirinlabs/utils/str"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"summersea.top/filetransfer/util"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

const ContentTypeJsonValue = "application/json;charset=UTF-8"

type FileServer struct {
	http.Handler
	store Store
}

func NewFileServer(store Store) *FileServer {
	fileServer := &FileServer{}
	router := http.NewServeMux()
	router.Handle("/file/upload/initialization", http.HandlerFunc(fileServer.uploadInitHandler))
	router.Handle("/file/upload", http.HandlerFunc(fileServer.uploadHandler))

	fileServer.Handler = router
	fileServer.store = store
	return fileServer
}

func (fs *FileServer) uploadInitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusForbidden)
		return
	}
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
	taskId := fs.handleUploadInit(*uploadInitBody)
	writeStringToResponse(w, taskId)
}

func (fs *FileServer) uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	param := util.ExtractUrlParam(r.URL.String())
	taskId := param["taskId"]
	uploadData := fs.store.GetUploadData(taskId)

	if uploadData == nil {
		taskNotFoundBody := ErrorBody{
			Error: ErrorContent{
				Message: "The task id is not found.",
				Code:    "ResourceNotFound",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-type", ContentTypeJsonValue)
		writeStructToResponse(w, taskNotFoundBody)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func (fs *FileServer) handleUploadInit(body UploadInitReqBody) string {
	return uuid.NewV4().String()
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

func writeStructToResponse(w http.ResponseWriter, object interface{}) {
	err := json.NewEncoder(w).Encode(object)
	if err != nil {
		log.Printf("problem encode struct %+v to json. err: %v", object, err)
	}
}

func writeStringToResponse(w http.ResponseWriter, data string) {
	_, err := w.Write([]byte(data))
	if err != nil {
		log.Printf("problem write data to response")
	}
}

type UploadData struct {
}

type Store interface {
	GetUploadData(taskId string) *UploadData
}
