package storage

import (
	"bufio"
	"os"
	"sync"
)

type FileStorage struct {
	filename     string
	calculations []string
	mutex        sync.RWMutex
}

func NewFileStorage(filename string) (*FileStorage, error) {
	storage := &FileStorage{
		filename:     filename,
		calculations: make([]string, 0),
	}

	// Load existing calculations on startup
	if err := storage.load(); err != nil {
		return nil, err
	}

	return storage, nil
}

func (f *FileStorage) Store(calc string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.calculations = append(f.calculations, calc)
	return
}

func (f *FileStorage) GetRecent(n int) []string {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	start := len(f.calculations) - n
	if start < 0 {
		start = 0
	}
	// we don't want to return the original slice to avoid external modification
	calcCopy := make([]string, len(f.calculations[start:]))
	for i := 0; i < len(calcCopy); i++ {
		calcCopy[i] = f.calculations[start+i]
	}
	return calcCopy
}

func (f *FileStorage) load() error {
	file, err := os.Open(f.filename)
	if os.IsNotExist(err) {
		// File doesn't exist yet - that's OK
		return nil
	}
	if err != nil {
		return err
	}
	defer file.Close()

	f.calculations = make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" { // Skip empty lines
			f.calculations = append(f.calculations, line)
		}
	}

	return scanner.Err()
}

func (f *FileStorage) save() error {
	file, err := os.OpenFile(f.filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, calc := range f.calculations {
		_, err := writer.WriteString(calc + "\n")
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

func (f *FileStorage) Close() error {
	return f.save() // Ensure data is saved when closing
}
