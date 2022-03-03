package filetransfer

import (
	"io"
)

type UploadData UploadInitReqBody

type FileTranDataAdapter struct {
	dataStore DataStore
}

func NewFileTranDataAdapter(store DataStore) *FileTranDataAdapter {
	return &FileTranDataAdapter{store}
}

func (f *FileTranDataAdapter) SaveUploadData(taskId string, uploadData UploadData) {
	f.dataStore.SaveUploadData(taskId, uploadData)
}

func (f *FileTranDataAdapter) IsTaskExist(taskId string) bool {
	return f.dataStore.IsTaskExist(taskId)
}

func (f *FileTranDataAdapter) GetUploadChannel(taskId string) io.WriteCloser {
	f.dataStore.GetUploadData(taskId)
	return nil
}

type DataStore interface {
	SaveUploadData(taskId string, data UploadData)
	GetUploadData(taskId string) *UploadData
	IsTaskExist(taskId string) bool
}
