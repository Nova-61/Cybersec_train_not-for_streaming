package storage

type MemoryStorage struct {
	data map[string][]byte
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string][]byte),
	}
}

func (m *MemoryStorage) Save(key string, data []byte) error {
	m.data[key] = data
	return nil
}

func (m *MemoryStorage) Load(key string) ([]byte, error) {
	data, ok := m.data[key]
	if !ok {
		return nil, nil
	}
	return data, nil
}
