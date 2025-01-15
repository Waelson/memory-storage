package resp

import (
	"fmt"
	"net"

	"github.com/Waelson/memory-storage/m-storage/internal/server"
)

// handleGet processa o comando GET.
func handleGet(args []string, conn net.Conn, storage *server.MStorageServer) {
	if len(args) != 2 {
		conn.Write([]byte("-ERR wrong number of arguments for 'GET' command\r\n"))
		return
	}
	key := args[1]
	value, found := storage.Get(key)
	if found {
		conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)))
	} else {
		conn.Write([]byte("$-1\r\n"))
	}
}
