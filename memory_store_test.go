package filetransfer_test

import (
	"summersea.top/filetransfer"
	"testing"
)

func TestMemoryStore_GetUploadData(t *testing.T) {
	t.Run("test get non exist data", func(t *testing.T) {
		store := filetransfer.NewMemoryStore()
		data := store.GetUploadDataWithRm(filetransfer.NewTaskId())
		assertNilUploadData(t, data)
	})

	t.Run("test rm data", func(t *testing.T) {
		store := filetransfer.NewMemoryStore()
		taskId := filetransfer.NewTaskId()
		store.SaveUploadData(taskId, filetransfer.UploadData{})
		data := store.GetUploadDataWithRm(taskId)
		assertNotNil(t, data)
		assertFalse(t, store.IsUploadTaskExist(taskId))
	})
}

func TestMemoryStore_IsTaskExist(t *testing.T) {
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
			got := m.GetUploadDataWithRm(tt.args.taskId)
			assertStructEquals(t, *got, tt.args.data)
		})
	}

	t.Run("test empty taskId", func(t *testing.T) {
		m := filetransfer.NewMemoryStore()
		saved := filetransfer.UploadData{}
		m.SaveUploadData("", saved)
		got := m.GetUploadDataWithRm("")
		assertNilUploadData(t, got)
	})
}

func assertNilUploadData(t *testing.T, got *filetransfer.UploadData) {
	t.Helper()
	if got != nil {
		t.Errorf("want nil but got other: %+v", got)
	}
}
