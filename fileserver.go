package filetransfer

import (
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
	r.POST("/file/download/initialization", fileServer.downloadInitHandler)
	r.GET("/file/download", fileServer.downloadHandler)
	fileServer.dataAdapter = adapter
	return r
}

func (fs *FileServerController) uploadInitHandler(ctx *gin.Context) {
	var uploadInitBody UploadInitReqBody
	if err := ctx.ShouldBindJSON(&uploadInitBody); err != nil {
		ctx.JSON(http.StatusBadRequest, getInvalidParamErr())
		return
	}
	if !fs.isUploadInitReqBodyValid(uploadInitBody) {
		ctx.JSON(http.StatusBadRequest, getInvalidParamErr())
		return
	}
	taskId := fs.handleUploadInit(UploadData(uploadInitBody))
	ctx.JSON(http.StatusOK, OkBody{Data: Data{"taskId": taskId}})
}

func (fs *FileServerController) handleUploadInit(uploadData UploadData) string {
	taskId := NewTaskId()
	fs.dataAdapter.SaveUploadData(taskId, uploadData)
	return taskId
}

func (fs *FileServerController) downloadInitHandler(ctx *gin.Context) {
	var downloadInitBody DownloadInitReqBody
	if err := ctx.ShouldBindJSON(&downloadInitBody); err != nil {
		ctx.JSON(http.StatusBadRequest, getInvalidParamErr())
		return
	}
	if !fs.isDownloadInitReqBodyValid(downloadInitBody) {
		ctx.JSON(http.StatusBadRequest, getInvalidParamErr())
		return
	}
	taskId := fs.handleDownloadInit(DownloadData(downloadInitBody))
	ctx.JSON(http.StatusOK, OkBody{Data: Data{"taskId": taskId}})
}

func (fs *FileServerController) handleDownloadInit(downloadData DownloadData) string {
	taskId := NewTaskId()
	fs.dataAdapter.SaveDownloadData(taskId, downloadData)
	return taskId
}

func (fs *FileServerController) uploadHandler(ctx *gin.Context) {
	taskId := ctx.Query("taskId")
	if !fs.dataAdapter.IsUploadTaskExist(taskId) {
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

func (fs *FileServerController) handleUpload(taskId string, reader io.Reader) error {
	writeCloser, err := fs.dataAdapter.GetUploadChannel(taskId)
	if err != nil {
		return fmt.Errorf("problem create upload channel %v", err)
	}
	defer closeWithErrLog(writeCloser)
	_, err = io.Copy(writeCloser, reader)
	if err != nil {
		return fmt.Errorf("problem transfer file: %v", err)
	}
	return nil
}

// 下载API的处理器，负责view部分的业务
func (fs *FileServerController) downloadHandler(ctx *gin.Context) {
	taskId := ctx.Query("taskId")
	setFilename := func(value string) {
		ctx.Writer.Header().Set("Content-Disposition", "attachment; filename="+value)
	}
	if !fs.dataAdapter.IsDownloadTaskExist(taskId) {
		ctx.JSON(http.StatusBadRequest, getTaskNotFoundErr())
	} else {
		err := fs.handleDownload(taskId, ctx.Writer, setFilename)
		if err != nil {
			log.Printf("problem download file: %v", err)
			ctx.Status(http.StatusInternalServerError)
		}
	}
}

func (fs *FileServerController) handleDownload(taskId string, writer io.Writer, setFilename func(value string)) error {
	readCloser, filename, err := fs.dataAdapter.GetDownloadChannelFilename(taskId)
	if err != nil {
		return fmt.Errorf("problem create download channel %v", err)
	}
	setFilename(filename)
	defer closeWithErrLog(readCloser)
	_, err = io.Copy(writer, readCloser)
	if err != nil {
		return fmt.Errorf("problem transfer file: %v", err)
	}
	return nil
}

func (fs *FileServerController) isUploadInitReqBodyValid(body UploadInitReqBody) bool {
	if !str.StartsWith(body.Path, "/") {
		return false
	}
	if body.Filename == "" || str.StartsWith(body.Filename, "/") {
		return false
	}
	return fs.isResourceReqBodyValid(body.Resource)
}

func (fs *FileServerController) isDownloadInitReqBodyValid(body DownloadInitReqBody) bool {
	if body.Path == "" || str.EndsWith(body.Path, "/") {
		return false
	}
	if !fs.isValidPathInLinux(body.Path) && !fs.isValidPathInWindows(body.Path) {
		return false
	}
	return fs.isResourceReqBodyValid(body.Resource)
}

func (fs *FileServerController) isValidPathInLinux(path string) bool {
	return str.StartsWith(path, "/")
}

func (fs *FileServerController) isValidPathInWindows(path string) bool {
	driveLetter := path[0]
	sep := path[1:3]
	return 64 < driveLetter && driveLetter < 91 && sep == ":\\"
}

func (fs *FileServerController) isResourceReqBodyValid(resource Resource) bool {
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

type DataAdapter interface {
	IsUploadTaskExist(taskId string) bool
	GetUploadChannel(taskId string) (WriteCloseRollback, error)
	SaveUploadData(taskId string, uploadData UploadData)
	IsDownloadTaskExist(taskId string) bool
	// GetDownloadChannelFilename 获取下载通道，并获取下载的文件名
	GetDownloadChannelFilename(taskId string) (io.ReadCloser, string, error)
	SaveDownloadData(taskId string, downloadData DownloadData)
}

func NewTaskId() string {
	return uuid.NewV4().String()
}
