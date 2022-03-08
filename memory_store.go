package filetransfer

type MemoryStore struct {
	uploadStore   map[string]UploadData
	downloadStore map[string]DownloadData
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		uploadStore:   make(map[string]UploadData),
		downloadStore: make(map[string]DownloadData),
	}
}

func (m *MemoryStore) SaveUploadData(taskId string, data UploadData) {
	if taskId == "" {
		return
	}
	m.uploadStore[taskId] = data
}

func (m *MemoryStore) GetUploadDataRemove(taskId string) *UploadData {
	if !m.IsUploadTaskExist(taskId) {
		return nil
	}
	data := m.uploadStore[taskId]
	m.removeTaskId(taskId)
	return &data
}

func (m *MemoryStore) IsUploadTaskExist(taskId string) bool {
	_, exist := m.uploadStore[taskId]
	return exist
}

func (m *MemoryStore) SaveDownloadData(taskId string, data DownloadData) {
	if taskId == "" {
		return
	}
	m.downloadStore[taskId] = data
}

func (m *MemoryStore) GetDownloadDataRemove(taskId string) *DownloadData {
	if !m.IsDownloadTaskExist(taskId) {
		return nil
	}
	data := m.downloadStore[taskId]
	m.removeTaskId(taskId)
	return &data
}

func (m *MemoryStore) IsDownloadTaskExist(taskId string) bool {
	_, exist := m.downloadStore[taskId]
	return exist
}

func (m *MemoryStore) removeTaskId(taskId string) {
	delete(m.uploadStore, taskId)
}
