package resp

import (
	"net"
	"strconv"
	"time"

	"github.com/Waelson/memory-storage/m-storage/internal/server"
)

// handleExpire processa o comando EXPIRE.
func handleExpire(args []string, conn net.Conn, storage *server.MStorageServer) {
	if len(args) != 3 {
		conn.Write([]byte("-ERR wrong number of arguments for 'EXPIRE' command\r\n"))
		return
	}

	key := args[1]
	seconds, err := strconv.Atoi(args[2])
	if err != nil || seconds < 0 {
		conn.Write([]byte("-ERR invalid TTL\r\n"))
		return
	}

	storage.LockData()
	defer storage.UnlockData()

	kv, exists := storage.GetData()[key]
	if exists {
		kv.Expiration = time.Now().Add(time.Duration(seconds) * time.Second)
		storage.GetData()[key] = kv
		conn.Write([]byte(":1\r\n")) // Sucesso
	} else {
		conn.Write([]byte(":0\r\n")) // Chave não existe
	}
}
