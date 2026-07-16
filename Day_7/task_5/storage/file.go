package storage

import (
	"os"
	"path/filepath"
)

type FileStorage struct {
	dir string
}

func NewFileStorage(dir string) (*FileStorage, error) {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}
	return &FileStorage{dir: dir}, nil
}

func (f *FileStorage) Save(key string, data []byte) error {
	path := filepath.Join(f.dir, key)
	return os.WriteFile(path, data, 0644)
}

func (f *FileStorage) Load(key string) ([]byte, error) {
	path := filepath.Join(f.dir, key)
	return os.ReadFile(path)
}
