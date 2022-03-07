package filetransfer_test

import (
	"log"
	"summersea.top/filetransfer"
	"testing"
)

type StubDataStore struct {
	saveCalls    int
	existCalls   int
	getDataCalls int
	taskId       string
	uploadData   filetransfer.UploadData
}

func (s *StubDataStore) SaveUploadData(taskId string, data filetransfer.UploadData) {
	s.saveCalls++
}

func (s *StubDataStore) GetUploadDataWithRm(taskId string) *filetransfer.UploadData {
	s.getDataCalls++
	if taskId == s.taskId {
		s.taskId = ""
		return &s.uploadData
	}
	return nil
}

func (s *StubDataStore) IsTaskExist(taskId string) bool {
	s.existCalls++
	return s.taskId == taskId
}

func TestFileTranDataAdapter_SaveUploadData(t *testing.T) {
	store := &StubDataStore{}
	adapter := filetransfer.NewFileTranDataAdapter(store)
	adapter.SaveUploadData("", filetransfer.UploadData{})
	assertIntEquals(t, store.saveCalls, 1)
}

func TestFileTranDataAdapter_IsTaskExist(t *testing.T) {
	existedTaskId := filetransfer.NewTaskId()
	missedTaskId := filetransfer.NewTaskId()
	store := &StubDataStore{taskId: existedTaskId}
	adapter := filetransfer.NewFileTranDataAdapter(store)
	assertTrue(t, adapter.IsUploadTaskExist(existedTaskId))
	assertFalse(t, adapter.IsUploadTaskExist(missedTaskId))
	assertIntEquals(t, store.existCalls, 2)
}

// 该测试需要配置外部sftp环境以测试，没有环境时可以无法通过
func TestFileTranDataAdapter_GetUploadData(t *testing.T) {
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
	assertIntEquals(t, store.getDataCalls, 1)
	if channel != nil {
		assertNil(t, channel.RollBack())
		assertNil(t, channel.Close())
	}
	assertFalse(t, adapter.IsUploadTaskExist(existedTaskId))
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
