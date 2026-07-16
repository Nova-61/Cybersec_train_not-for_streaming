package storage

type Storage interface {
	Save(key string, data []byte) error
	Load(key string) ([]byte, error)
}
