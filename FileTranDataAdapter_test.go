package filetransfer_test

import (
	"log"
	"summersea.top/filetransfer"
	"testing"
)

type StubDataStore struct {
	saveUploadCalls         int
	uploadExistCalls        int
	getUploadChannelCalls   int
	taskId                  string
	saveDownloadCalls       int
	downloadExistCalls      int
	getDownloadChannelCalls int
	uploadData              filetransfer.UploadData
	downloadData            filetransfer.DownloadData
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
	s.saveDownloadCalls++
}

func (s *StubDataStore) GetDownloadDataRemove(taskId string) *filetransfer.DownloadData {
	s.getDownloadChannelCalls++
	if taskId == s.taskId {
		s.taskId = ""
		return &s.downloadData
	}
	return nil
}

func (s *StubDataStore) IsDownloadTaskExist(taskId string) bool {
	s.downloadExistCalls++
	return s.taskId == taskId
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
		Resource: getSftpResource(),
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
	store := &StubDataStore{}
	adapter := filetransfer.NewFileTranDataAdapter(store)
	adapter.SaveDownloadData(filetransfer.NewTaskId(), filetransfer.DownloadData{})
	assertIntEquals(t, store.saveDownloadCalls, 1)
}

func TestFileTranDataAdapter_IsDownloadTaskExist(t *testing.T) {
	existedTaskId := filetransfer.NewTaskId()
	missedTaskId := filetransfer.NewTaskId()
	store := &StubDataStore{taskId: existedTaskId}
	adapter := filetransfer.NewFileTranDataAdapter(store)
	assertTrue(t, adapter.IsDownloadTaskExist(existedTaskId))
	assertFalse(t, adapter.IsDownloadTaskExist(missedTaskId))
	assertIntEquals(t, store.downloadExistCalls, 2)
}

// 测试点
// Path指定的目标为目录
// 析出filename
// sftp连接成功
// 任务号清空
func TestFileTranDataAdapter_GetDownloadChannelFilename(t *testing.T) {
	existedTaskId := filetransfer.NewTaskId()
	t.Run("common test", func(t *testing.T) {
		store := &StubDataStore{taskId: existedTaskId, downloadData: filetransfer.DownloadData{
			Resource: getSftpResource(),
			// 需要目标机器有该文件
			Path: "/home/test/ccc.txt",
		}}
		adapter := filetransfer.NewFileTranDataAdapter(store)
		channel, filename, err := adapter.GetDownloadChannelFilename(existedTaskId)
		if err != nil {
			log.Printf("%v", err)
		}
		assertNotNil(t, channel)
		assertIntEquals(t, store.getDownloadChannelCalls, 1)
		if channel != nil {
			assertNil(t, channel.Close())
		}
		assertStringEqual(t, filename, "ccc.txt")
		assertFalse(t, adapter.IsUploadTaskExist(existedTaskId))
	})

	t.Run("input path without filename", func(t *testing.T) {
		store := &StubDataStore{taskId: existedTaskId, downloadData: filetransfer.DownloadData{
			Resource: getSftpResource(),
			Path:     "/home/test",
		}}
		adapter := filetransfer.NewFileTranDataAdapter(store)
		_, _, err := adapter.GetDownloadChannelFilename(existedTaskId)
		assertErrEquals(t, err, filetransfer.DownloadDir)
	})
}

func getSftpResource() filetransfer.Resource {
	return filetransfer.Resource{
		Address: "localhost",
		Port:    22,
		Account: filetransfer.Account{
			Name:     "test",
			Password: "test",
		},
	}
}
