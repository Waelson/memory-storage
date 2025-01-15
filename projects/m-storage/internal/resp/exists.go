package resp

import (
	"net"

	"github.com/Waelson/memory-storage/m-storage/internal/server"
)

// handleExists processa o comando EXISTS.
func handleExists(args []string, conn net.Conn, storage *server.MStorageServer) {
	if len(args) != 2 {
		conn.Write([]byte("-ERR wrong number of arguments for 'EXISTS' command\r\n"))
		return
	}
	key := args[1]
	_, exists := storage.Get(key)
	if exists {
		conn.Write([]byte(":1\r\n"))
	} else {
		conn.Write([]byte(":0\r\n"))
	}
}
