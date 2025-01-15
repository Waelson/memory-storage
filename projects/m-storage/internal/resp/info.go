package resp

import (
	"fmt"
	"net"
	"time"

	"github.com/Waelson/memory-storage/m-storage/internal/server"
)

// handleInfo processa o comando INFO.
func handleInfo(args []string, conn net.Conn, storage *server.MStorageServer) {
	if len(args) != 1 {
		conn.Write([]byte("-ERR wrong number of arguments for 'INFO' command\r\n"))
		return
	}

	info := generateInfo(storage)
	conn.Write([]byte(info))
}

// generateInfo gera estatísticas sobre o estado do servidor.
func generateInfo(storage *server.MStorageServer) string {
	// Tempo de execução
	uptime := time.Since(storage.GetStartTime()).Seconds()

	// Protege o acesso ao campo data com um lock
	storage.LockData()
	defer storage.UnlockData()

	// Número de chaves armazenadas
	numKeys := len(storage.GetData())

	info := fmt.Sprintf(
		"# Server\r\nuptime_in_seconds:%.0f\r\n\r\n"+
			"# Stats\r\nnumber_of_keys:%d\r\n\r\n"+
			"# Memory\r\nused_memory:%d\r\n",
		uptime,
		numKeys,
		estimateMemoryUsage(storage.GetData()), // Estima a memória usada
	)

	return info
}

// estimateMemoryUsage estima o uso de memória do servidor.
func estimateMemoryUsage(data map[string]server.KeyValue) int {
	// Estimativa simples para cada chave e valor
	const avgSizePerEntry = 100 // Valor médio em bytes por chave/valor
	return len(data) * avgSizePerEntry
}
