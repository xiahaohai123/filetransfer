package filetransfer_test

import (
	"summersea.top/filetransfer"
	"testing"
)

func TestMemoryStore_GetUploadDataRemove(t *testing.T) {
	t.Run("test get non exist data", func(t *testing.T) {
		store := filetransfer.NewMemoryStore()
		data := store.GetUploadDataRemove(filetransfer.NewTaskId())
		assertNil(t, data)
	})

	t.Run("test rm data", func(t *testing.T) {
		store := filetransfer.NewMemoryStore()
		taskId := filetransfer.NewTaskId()
		store.SaveUploadData(taskId, filetransfer.UploadData{})
		data := store.GetUploadDataRemove(taskId)
		assertNotNil(t, data)
		assertFalse(t, store.IsUploadTaskExist(taskId))
	})
}

func TestMemoryStore_IsUploadTaskExist(t *testing.T) {
	t.Run("test exist data", func(t *testing.T) {
		store := filetransfer.NewMemoryStore()
		taskId := filetransfer.NewTaskId()
		store.SaveUploadData(taskId, filetransfer.UploadData{})
		assertTrue(t, store.IsUploadTaskExist(taskId))
	})

	t.Run("test non exist data", func(t *testing.T) {
		store := filetransfer.NewMemoryStore()
		assertFalse(t, store.IsUploadTaskExist(filetransfer.NewTaskId()))

		store.SaveUploadData("", filetransfer.UploadData{})
		assertFalse(t, store.IsUploadTaskExist(""))
	})
}

func TestMemoryStore_SaveUploadData(t *testing.T) {
	type argsAndWant struct {
		taskId string
		data   filetransfer.UploadData
	}
	testCases := []argsAndWant{
		{filetransfer.NewTaskId(), filetransfer.UploadData{}},
		{filetransfer.NewTaskId(), filetransfer.UploadData{
			Resource: filetransfer.Resource{Address: "a", Port: 22,
				Account: filetransfer.Account{Name: "a", Password: "a"}},
			Filename: "aaa", Path: "aaa"},
		},
	}
	tests := []struct {
		name string
		args argsAndWant
	}{
		{name: "test01", args: testCases[0]},
		{name: "test02", args: testCases[1]},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := filetransfer.NewMemoryStore()
			m.SaveUploadData(tt.args.taskId, tt.args.data)
			got := m.GetUploadDataRemove(tt.args.taskId)
			assertStructEquals(t, *got, tt.args.data)
		})
	}

	t.Run("test empty taskId", func(t *testing.T) {
		m := filetransfer.NewMemoryStore()
		saved := filetransfer.UploadData{}
		m.SaveUploadData("", saved)
		got := m.GetUploadDataRemove("")
		assertNil(t, got)
	})
}

func TestMemoryStore_GetDownloadDataRemove(t *testing.T) {
	t.Run("test get non exist data", func(t *testing.T) {
		store := filetransfer.NewMemoryStore()
		data := store.GetDownloadDataRemove(filetransfer.NewTaskId())
		assertNil(t, data)
	})

	t.Run("test rm data", func(t *testing.T) {
		store := filetransfer.NewMemoryStore()
		taskId := filetransfer.NewTaskId()
		store.SaveDownloadData(taskId, filetransfer.DownloadData{})
		data := store.GetDownloadDataRemove(taskId)
		assertNotNil(t, data)
		assertFalse(t, store.IsUploadTaskExist(taskId))
	})
}

func TestMemoryStore_IsDownloadTaskExist(t *testing.T) {
	t.Run("test exist data", func(t *testing.T) {
		store := filetransfer.NewMemoryStore()
		taskId := filetransfer.NewTaskId()
		store.SaveDownloadData(taskId, filetransfer.DownloadData{})
		assertTrue(t, store.IsDownloadTaskExist(taskId))
	})

	t.Run("test non exist data", func(t *testing.T) {
		store := filetransfer.NewMemoryStore()
		assertFalse(t, store.IsDownloadTaskExist(filetransfer.NewTaskId()))

		store.SaveDownloadData("", filetransfer.DownloadData{})
		assertFalse(t, store.IsDownloadTaskExist(""))
	})
}

func TestMemoryStore_SaveDownloadData(t *testing.T) {
	type argsAndWant struct {
		taskId string
		data   filetransfer.DownloadData
	}
	testCases := []argsAndWant{
		{filetransfer.NewTaskId(), filetransfer.DownloadData{}},
		{filetransfer.NewTaskId(), filetransfer.DownloadData{
			Resource: filetransfer.Resource{Address: "a", Port: 22,
				Account: filetransfer.Account{Name: "a", Password: "a"}}},
		},
	}
	tests := []struct {
		name string
		args argsAndWant
	}{
		{name: "test01", args: testCases[0]},
		{name: "test02", args: testCases[1]},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := filetransfer.NewMemoryStore()
			m.SaveDownloadData(tt.args.taskId, tt.args.data)
			got := m.GetDownloadDataRemove(tt.args.taskId)
			assertStructEquals(t, *got, tt.args.data)
		})
	}

	t.Run("test empty taskId", func(t *testing.T) {
		m := filetransfer.NewMemoryStore()
		saved := filetransfer.DownloadData{}
		m.SaveDownloadData("", saved)
		got := m.GetDownloadDataRemove("")
		assertNil(t, got)
	})
}
