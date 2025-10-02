package storage

type Storage interface {
	Store(string)
	GetRecent(int) []string
	save() error
	load() error
	Close() error
}

func NewStorage(storageType, filename string) (Storage, error) {
	// creating const for single use seems unnecessary
	switch storageType {
	case "file":
		return NewFileStorage(filename)
	case "memory":
		return NewMemoryStorage(), nil
	default:
		return nil, nil
	}
}
