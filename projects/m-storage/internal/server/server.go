package server

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"sync"
	"time"
)

// KeyValue representa uma entrada com valor e tempo de expiração.
type KeyValue struct {
	Value      string
	Expiration time.Time
}

// MStorageServer é a estrutura do servidor Redis-like.
type MStorageServer struct {
	data      map[string]KeyValue // Armazena os dados
	mu        sync.RWMutex        // Protege o acesso aos dados
	filePath  string              // Caminho do arquivo de persistência
	writeCh   chan struct{}       // Canal para escrita assíncrona
	startTime time.Time           // Tempo de início do servidor
}

// NewRedisServer cria uma nova instância do servidor.
func NewRedisServer(filePath string) *MStorageServer {
	server := &MStorageServer{
		data:     make(map[string]KeyValue),
		filePath: filePath,
		writeCh:  make(chan struct{}, 1), // Canal para controle de escrita assíncrona.
	}
	server.loadFromDisk()
	go server.asyncWriteToDisk()
	go server.cleanupExpiredKeys()
	return server
}

// Set define uma chave com valor e tempo de expiração (em segundos).
func (r *MStorageServer) Set(key, value string, ttl time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	expiration := time.Now().Add(ttl)
	r.data[key] = KeyValue{
		Value:      value,
		Expiration: expiration,
	}
	select {
	case r.writeCh <- struct{}{}:
	default:
	}
}

// Get retorna o valor de uma chave, se não estiver expirada.
func (r *MStorageServer) Get(key string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if kv, exists := r.data[key]; exists {
		if kv.Expiration.After(time.Now()) {
			return kv.Value, true
		}
		delete(r.data, key)
	}
	return "", false
}

// Del remove uma chave do servidor.
func (r *MStorageServer) Del(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.data[key]; exists {
		delete(r.data, key)
		select {
		case r.writeCh <- struct{}{}:
		default:
		}
		return true
	}
	return false
}

// cleanupExpiredKeys remove chaves expiradas periodicamente.
func (r *MStorageServer) cleanupExpiredKeys() {
	for {
		time.Sleep(1 * time.Second)
		r.mu.Lock()
		for key, kv := range r.data {
			if kv.Expiration.Before(time.Now()) {
				delete(r.data, key)
			}
		}
		r.mu.Unlock()
	}
}

// loadFromDisk carrega os dados do disco para a memória.
func (r *MStorageServer) loadFromDisk() {
	file, err := os.Open(r.filePath)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println("Error loading data from disk:", err)
		}
		return
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&r.data); err != nil {
		fmt.Println("Error decoding data from disk:", err)
	}
	fmt.Println("Data loaded from disk.")
}

// asyncWriteToDisk grava os dados em disco de forma assíncrona.
func (r *MStorageServer) asyncWriteToDisk() {
	for range r.writeCh {
		r.mu.RLock()
		data := r.data
		r.mu.RUnlock()
		var buffer bytes.Buffer
		encoder := gob.NewEncoder(&buffer)
		if err := encoder.Encode(data); err != nil {
			fmt.Println("Error encoding data to disk:", err)
			continue
		}
		file, err := os.Create(r.filePath)
		if err != nil {
			fmt.Println("Error writing data to disk:", err)
			continue
		}
		_, err = file.Write(buffer.Bytes())
		if err != nil {
			fmt.Println("Error writing binary data to disk:", err)
		}
		file.Close()
	}
}

// Métodos de acesso ao campo startTime
func (s *MStorageServer) GetStartTime() time.Time {
	return s.startTime
}

// Métodos de acesso ao campo data
func (s *MStorageServer) GetData() map[string]KeyValue {
	return s.data
}

// LockData e UnlockData para manipular o campo data de forma segura
func (s *MStorageServer) LockData() {
	s.mu.RLock()
}

func (s *MStorageServer) UnlockData() {
	s.mu.RUnlock()
}

// SignalWrite sinaliza que novos dados precisam ser escritos no disco.
func (s *MStorageServer) SignalWrite(data map[string]KeyValue) {
	select {
	case s.writeCh <- struct{}{}:
	default:
		// Evita múltiplas sinalizações redundantes
	}
}
