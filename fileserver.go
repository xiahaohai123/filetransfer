package filetransfer

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kirinlabs/utils/str"
	"github.com/satori/go.uuid"
	"io"
	"log"
	"net/http"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

type FileServerController struct {
	dataAdapter DataAdapter
}

func NewFileServer(adapter DataAdapter) *gin.Engine {
	fileServer := &FileServerController{}
	r := gin.Default()
	r.POST("/file/upload/initialization", fileServer.uploadInitHandler)
	r.POST("/file/upload", fileServer.uploadHandler)
	fileServer.dataAdapter = adapter
	return r
}

func (fs *FileServerController) uploadInitHandler(ctx *gin.Context) {
	var uploadInitBody UploadInitReqBody
	if err := ctx.ShouldBindJSON(&uploadInitBody); err != nil {
		ctx.JSON(http.StatusBadRequest, getInvalidParamErr())
		return
	}
	if !isUploadInitReqBodyValid(uploadInitBody) {
		ctx.JSON(http.StatusBadRequest, getInvalidParamErr())
		return
	}
	taskId := fs.handleUploadInit(UploadData(uploadInitBody))
	ctx.String(http.StatusOK, taskId)
}

func (fs *FileServerController) uploadHandler(ctx *gin.Context) {
	taskId := ctx.Query("taskId")
	if !fs.dataAdapter.IsTaskExist(taskId) {
		ctx.JSON(http.StatusBadRequest, getTaskNotFoundErr())
	} else {
		err := fs.handleUpload(taskId, ctx.Request.Body)
		if err != nil {
			log.Printf("problem upload file: %v", err)
			ctx.Status(http.StatusInternalServerError)
		} else {
			ctx.Status(http.StatusNoContent)
		}
	}
}

func (fs *FileServerController) handleUploadInit(uploadData UploadData) string {
	taskId := NewTaskId()
	fs.dataAdapter.SaveUploadData(taskId, uploadData)
	return taskId
}

func (fs *FileServerController) extractBody(r *http.Request) (*UploadInitReqBody, error) {
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

func (fs *FileServerController) handleUpload(taskId string, reader io.Reader) error {
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

type DataAdapter interface {
	IsTaskExist(taskId string) bool
	GetUploadChannel(taskId string) (WriteCloseRollback, error)
	SaveUploadData(taskId string, uploadData UploadData)
}

func NewTaskId() string {
	return uuid.NewV4().String()
}
