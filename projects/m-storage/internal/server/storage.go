package server

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"sync"
	"syscall"
)

// StorageHandler gerencia a persistência de dados no disco usando mmap.
type StorageHandler struct {
	filePath string
	writeCh  chan struct{}
	mu       sync.Mutex
}

// NewStorageHandler cria uma nova instância do StorageHandler.
func NewStorageHandler(filePath string) *StorageHandler {
	handler := &StorageHandler{
		filePath: filePath,
		writeCh:  make(chan struct{}, 1),
	}

	go handler.asyncWriteLoop()
	return handler
}

// Load carrega os dados do arquivo mmap para a memória.
func (s *StorageHandler) Load() map[string]KeyValue {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.OpenFile(s.filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return make(map[string]KeyValue)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file stats:", err)
		return make(map[string]KeyValue)
	}

	if stat.Size() == 0 {
		return make(map[string]KeyValue)
	}

	// Mapeia o arquivo para memória
	data, err := syscall.Mmap(int(file.Fd()), 0, int(stat.Size()), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		fmt.Println("Error mapping file:", err)
		return make(map[string]KeyValue)
	}
	defer syscall.Munmap(data)

	// Decodifica os dados
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	loadedData := make(map[string]KeyValue)
	if err := decoder.Decode(&loadedData); err != nil {
		fmt.Println("Error decoding data:", err)
		return make(map[string]KeyValue)
	}

	fmt.Println("Data loaded from disk.")
	return loadedData
}

// SignalWrite sinaliza que novos dados precisam ser escritos no disco.
func (s *StorageHandler) SignalWrite(data map[string]KeyValue) {
	select {
	case s.writeCh <- struct{}{}:
	default:
		// Evita múltiplas sinalizações redundantes
	}
}

// asyncWriteLoop gerencia a escrita assíncrona dos dados no disco.
func (s *StorageHandler) asyncWriteLoop() {
	for range s.writeCh {
		s.mu.Lock()
		data := <-s.writeCh

		file, err := os.OpenFile(s.filePath, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			fmt.Println("Error opening file for writing:", err)
			s.mu.Unlock()
			continue
		}
		defer file.Close()

		var buffer bytes.Buffer
		encoder := gob.NewEncoder(&buffer)
		if err := encoder.Encode(data); err != nil {
			fmt.Println("Error encoding data:", err)
			s.mu.Unlock()
			continue
		}

		// Mapeia o arquivo para a memória
		fileSize := buffer.Len()
		if err := file.Truncate(int64(fileSize)); err != nil {
			fmt.Println("Error truncating file:", err)
			s.mu.Unlock()
			continue
		}

		mmap, err := syscall.Mmap(int(file.Fd()), 0, fileSize, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
		if err != nil {
			fmt.Println("Error mapping file for writing:", err)
			s.mu.Unlock()
			continue
		}
		copy(mmap, buffer.Bytes())
		syscall.Munmap(mmap)

		fmt.Println("Data successfully written to disk.")
		s.mu.Unlock()
	}
}
