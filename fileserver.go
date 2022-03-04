package filetransfer

import (
	"encoding/json"
	"fmt"
	"github.com/kirinlabs/utils/str"
	"github.com/satori/go.uuid"
	"io"
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
	dataAdapter DataAdapter
}

func NewFileServer(adapter DataAdapter) *FileServer {
	fileServer := &FileServer{}
	router := http.NewServeMux()
	router.Handle("/file/upload/initialization", http.HandlerFunc(fileServer.uploadInitHandler))
	router.Handle("/file/upload", http.HandlerFunc(fileServer.uploadHandler))

	fileServer.Handler = router
	fileServer.dataAdapter = adapter
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	taskId := fs.handleUploadInit(UploadData(*uploadInitBody))
	writeStringToResponse(w, taskId)
}

func (fs *FileServer) uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	param := util.ExtractUrlParam(r.URL.String())
	taskId := param["taskId"]

	if !fs.dataAdapter.IsTaskExist(taskId) {
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
		err := fs.handleUpload(taskId, r.Body)
		if err != nil {
			log.Printf("problem upload file: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	}
}

func (fs *FileServer) handleUploadInit(uploadData UploadData) string {
	taskId := NewTaskId()
	fs.dataAdapter.SaveUploadData(taskId, uploadData)
	return taskId
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
	if body.Filename == "" || str.StartsWith(body.Filename, "/") {
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

func (fs *FileServer) handleUpload(taskId string, reader io.Reader) error {
	writeCloser, err := fs.dataAdapter.GetUploadChannel(taskId)
	if err != nil {
		return fmt.Errorf("problem create channel %v", err)
	}
	defer closeWithErrLog(writeCloser)
	_, err = io.Copy(writeCloser, reader)
	if err != nil {
		return fmt.Errorf("problem transfer file: %v", err)
	}
	return nil
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

type DataAdapter interface {
	IsTaskExist(taskId string) bool
	GetUploadChannel(taskId string) (WriteCloseRollback, error)
	SaveUploadData(taskId string, uploadData UploadData)
}

func NewTaskId() string {
	return uuid.NewV4().String()
}
