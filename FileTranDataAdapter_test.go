package filetransfer_test

import (
	"log"
	"summersea.top/filetransfer"
	"testing"
)

type StubDataStore struct {
	saveUploadCalls       int
	uploadExistCalls      int
	getUploadChannelCalls int
	taskId                string
	uploadData            filetransfer.UploadData
	downloadData          filetransfer.DownloadData
}

func (s *StubDataStore) SaveUploadData(taskId string, data filetransfer.UploadData) {
	s.saveUploadCalls++
}

func (s *StubDataStore) GetUploadDataRemove(taskId string) *filetransfer.UploadData {
	s.getUploadChannelCalls++
	if taskId == s.taskId {
		s.taskId = ""
		return &s.uploadData
	}
	return nil
}

func (s *StubDataStore) IsUploadTaskExist(taskId string) bool {
	s.uploadExistCalls++
	return s.taskId == taskId
}

func (s *StubDataStore) SaveDownloadData(taskId string, data filetransfer.DownloadData) {
	panic("implement me")
}

func (s *StubDataStore) GetDownloadDataRemove(taskId string) *filetransfer.DownloadData {
	panic("implement me")
}

func (s *StubDataStore) IsDownloadTaskExist(taskId string) bool {
	panic("implement me")
}

func TestFileTranDataAdapter_SaveUploadData(t *testing.T) {
	store := &StubDataStore{}
	adapter := filetransfer.NewFileTranDataAdapter(store)
	adapter.SaveUploadData("", filetransfer.UploadData{})
	assertIntEquals(t, store.saveUploadCalls, 1)
}

func TestFileTranDataAdapter_IsTaskExist(t *testing.T) {
	existedTaskId := filetransfer.NewTaskId()
	missedTaskId := filetransfer.NewTaskId()
	store := &StubDataStore{taskId: existedTaskId}
	adapter := filetransfer.NewFileTranDataAdapter(store)
	assertTrue(t, adapter.IsUploadTaskExist(existedTaskId))
	assertFalse(t, adapter.IsUploadTaskExist(missedTaskId))
	assertIntEquals(t, store.uploadExistCalls, 2)
}

// 该测试需要配置外部sftp环境以测试，没有环境时可以无法通过
func TestFileTranDataAdapter_GetUploadChannel(t *testing.T) {
	existedTaskId := filetransfer.NewTaskId()
	store := &StubDataStore{taskId: existedTaskId, uploadData: filetransfer.UploadData{
		Resource: filetransfer.Resource{
			Address: "192.168.138.129",
			Port:    22,
			Account: filetransfer.Account{
				Name:     "test",
				Password: "test",
			},
		},
		Path:     "/home/test",
		Filename: "testAaa.txt",
	}}
	adapter := filetransfer.NewFileTranDataAdapter(store)
	channel, err := adapter.GetUploadChannel(existedTaskId)
	if err != nil {
		log.Printf("%v", err)
	}
	assertNotNil(t, channel)
	assertIntEquals(t, store.getUploadChannelCalls, 1)
	if channel != nil {
		assertNil(t, channel.RollBack())
		assertNil(t, channel.Close())
	}
	assertFalse(t, adapter.IsUploadTaskExist(existedTaskId))
}

func TestFileTranDataAdapter_SaveDownloadData(t *testing.T) {
}

func assertTrue(t *testing.T, got bool) {
	t.Helper()
	if !got {
		t.Errorf("want true but got false")
	}
}

func assertFalse(t *testing.T, got bool) {
	t.Helper()
	if got {
		t.Errorf("want false but got true")
	}
}

func assertNotNil(t *testing.T, got interface{}) {
	t.Helper()
	if got == nil {
		t.Errorf("want not nil but got")
	}
}

func assertNil(t *testing.T, got interface{}) {
	t.Helper()
	if got != nil {
		t.Errorf("want nil but got other: %+v", got)
	}
}
