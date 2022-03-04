package filetransfer

type MemoryStore struct {
	storeData map[string]UploadData
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		storeData: make(map[string]UploadData),
	}
}

func (m *MemoryStore) SaveUploadData(taskId string, data UploadData) {
	if taskId == "" {
		return
	}
	m.storeData[taskId] = data
}

func (m *MemoryStore) GetUploadData(taskId string) *UploadData {
	if !m.IsTaskExist(taskId) {
		return nil
	}
	data := m.storeData[taskId]
	return &data
}

func (m *MemoryStore) IsTaskExist(taskId string) bool {
	_, exist := m.storeData[taskId]
	return exist
}
