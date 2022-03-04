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

func (s *StubDataStore) GetUploadData(taskId string) *filetransfer.UploadData {
	s.getDataCalls++
	return &s.uploadData
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
	assertTrue(t, adapter.IsTaskExist(existedTaskId))
	assertFalse(t, adapter.IsTaskExist(missedTaskId))
	assertIntEquals(t, store.existCalls, 2)
}

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
}

func assertTrue(t *testing.T, got bool) {
	t.Helper()
	if !got {
		t.Errorf("want true bug got false")
	}
}

func assertFalse(t *testing.T, got bool) {
	t.Helper()
	if got {
		t.Errorf("want false bug got true")
	}
}

func assertNotNil(t *testing.T, got interface{}) {
	t.Helper()
	if got == nil {
		t.Errorf("want not nil bug got")
	}
}

func assertNil(t *testing.T, got interface{}) {
	t.Helper()
	if got != nil {
		t.Errorf("want nil but got other: %+v", got)
	}
}
