package resp

import (
	"fmt"
	"net"
	"time"

	"github.com/Waelson/memory-storage/m-storage/internal/server"
)

// handleTTL processa o comando TTL.
func handleTTL(args []string, conn net.Conn, storage *server.MStorageServer) {
	if len(args) != 2 {
		conn.Write([]byte("-ERR wrong number of arguments for 'TTL' command\r\n"))
		return
	}

	key := args[1]

	storage.LockData()
	defer storage.UnlockData()

	kv, exists := storage.GetData()[key]
	if !exists {
		conn.Write([]byte(":-2\r\n")) // Chave não existe
		return
	}

	if kv.Expiration.IsZero() {
		conn.Write([]byte(":-1\r\n")) // Chave não tem TTL
		return
	}

	ttl := int(time.Until(kv.Expiration).Seconds())
	if ttl <= 0 {
		conn.Write([]byte(":-2\r\n"))  // Chave expirada
		delete(storage.GetData(), key) // Remove chave expirada
	} else {
		conn.Write([]byte(fmt.Sprintf(":%d\r\n", ttl)))
	}
}
