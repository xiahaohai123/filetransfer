package filetransfer_test

import (
	"summersea.top/filetransfer"
	"testing"
)

func TestMemoryStore_GetUploadDataRemove(t *testing.T) {
	dataStores := createStores(t)

	for _, store := range dataStores {
		t.Run("test get non exist data", func(t *testing.T) {
			data := store.GetUploadDataRemove(filetransfer.NewTaskId())
			assertNil(t, data)
		})

		t.Run("test rm data", func(t *testing.T) {
			taskId := filetransfer.NewTaskId()
			store.SaveUploadData(taskId, filetransfer.UploadData{})
			data := store.GetUploadDataRemove(taskId)
			assertNotNil(t, data)
			assertFalse(t, store.IsUploadTaskExist(taskId))
		})
	}
}

func TestMemoryStore_IsUploadTaskExist(t *testing.T) {
	dataStores := createStores(t)
	for _, store := range dataStores {
		t.Run("test exist data", func(t *testing.T) {
			taskId := filetransfer.NewTaskId()
			store.SaveUploadData(taskId, filetransfer.UploadData{})
			assertTrue(t, store.IsUploadTaskExist(taskId))
		})

		t.Run("test non exist data", func(t *testing.T) {
			assertFalse(t, store.IsUploadTaskExist(filetransfer.NewTaskId()))

			store.SaveUploadData("", filetransfer.UploadData{})
			assertFalse(t, store.IsUploadTaskExist(""))
		})
	}
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
	dataStores := createStores(t)

	for _, store := range dataStores {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				store.SaveUploadData(tt.args.taskId, tt.args.data)
				got := store.GetUploadDataRemove(tt.args.taskId)
				assertStructEquals(t, *got, tt.args.data)
			})
		}

		t.Run("test empty taskId", func(t *testing.T) {
			saved := filetransfer.UploadData{}
			store.SaveUploadData("", saved)
			got := store.GetUploadDataRemove("")
			assertNil(t, got)
		})
	}
}

func TestMemoryStore_GetDownloadDataRemove(t *testing.T) {
	dataStores := createStores(t)
	for _, store := range dataStores {
		t.Run("test get non exist data", func(t *testing.T) {
			data := store.GetDownloadDataRemove(filetransfer.NewTaskId())
			assertNil(t, data)
		})

		t.Run("test rm data", func(t *testing.T) {
			taskId := filetransfer.NewTaskId()
			store.SaveDownloadData(taskId, filetransfer.DownloadData{})
			data := store.GetDownloadDataRemove(taskId)
			assertNotNil(t, data)
			assertFalse(t, store.IsUploadTaskExist(taskId))
		})
	}
}

func TestMemoryStore_IsDownloadTaskExist(t *testing.T) {
	dataStores := createStores(t)
	for _, store := range dataStores {
		t.Run("test exist data", func(t *testing.T) {
			taskId := filetransfer.NewTaskId()
			store.SaveDownloadData(taskId, filetransfer.DownloadData{})
			assertTrue(t, store.IsDownloadTaskExist(taskId))
		})

		t.Run("test non exist data", func(t *testing.T) {
			assertFalse(t, store.IsDownloadTaskExist(filetransfer.NewTaskId()))

			store.SaveDownloadData("", filetransfer.DownloadData{})
			assertFalse(t, store.IsDownloadTaskExist(""))
		})
	}
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
	dataStores := createStores(t)
	for _, store := range dataStores {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				store.SaveDownloadData(tt.args.taskId, tt.args.data)
				got := store.GetDownloadDataRemove(tt.args.taskId)
				assertStructEquals(t, *got, tt.args.data)
			})
		}

		t.Run("test empty taskId", func(t *testing.T) {
			saved := filetransfer.DownloadData{}
			store.SaveDownloadData("", saved)
			got := store.GetDownloadDataRemove("")
			assertNil(t, got)
		})
	}
}

func createStores(t *testing.T) []filetransfer.DataStore {
	redisStore, err := filetransfer.NewRedisStore("localhost:6379", "", 0)
	memoryStore := filetransfer.NewMemoryStore()
	dataStores := []filetransfer.DataStore{memoryStore}
	if err == nil {
		dataStores = append(dataStores, redisStore)
	} else {
		t.Errorf("problem create redis store: %v", err)
	}
	return dataStores
}
