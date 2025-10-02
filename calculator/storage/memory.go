package storage

import "sync"

type MemoryStorage struct {
	calculations []string
	mutex        sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		calculations: make([]string, 0),
	}
}

func (m *MemoryStorage) Store(calc string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.calculations = append(m.calculations, calc)
	return
}

func (m *MemoryStorage) GetRecent(n int) []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	start := len(m.calculations) - n
	if start < 0 {
		start = 0
	}
	// we don't want to return the original slice to avoid external modification
	calcCopy := make([]string, len(m.calculations[start:]))
	for i := 0; i < len(calcCopy); i++ {
		calcCopy[i] = m.calculations[start+i]
	}
	return calcCopy
}

func (m *MemoryStorage) save() error  { return nil } // Nothing to save for memory storage
func (m *MemoryStorage) load() error  { return nil } // Nothing to load for memory storage
func (m *MemoryStorage) Close() error { return nil } // Nothing to close for memory storage
