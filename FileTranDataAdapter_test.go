package filetransfer_test

import (
	"summersea.top/filetransfer"
	"testing"
)

type StubDataStore struct {
	saveCalls    int
	existCalls   int
	getDataCalls int
	taskId       string
}

func (s *StubDataStore) SaveUploadData(taskId string, data filetransfer.UploadData) {
	s.saveCalls++
}

func (s *StubDataStore) GetUploadData(taskId string) *filetransfer.UploadData {
	s.getDataCalls++
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
	assertTrue(t, adapter.IsTaskExist(existedTaskId))
	assertFalse(t, adapter.IsTaskExist(missedTaskId))
	assertIntEquals(t, store.existCalls, 2)
}

func TestFileTranDataAdapter_GetUploadData(t *testing.T) {
	existedTaskId := filetransfer.NewTaskId()
	store := &StubDataStore{taskId: existedTaskId}
	adapter := filetransfer.NewFileTranDataAdapter(store)
	channel := adapter.GetUploadChannel(existedTaskId)
	assertNotNil(t, channel)
	assertIntEquals(t, store.getDataCalls, 1)
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
