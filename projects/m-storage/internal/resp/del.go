package resp

import (
	"net"

	"github.com/Waelson/memory-storage/m-storage/internal/server"
)

// handleDel processa o comando DEL.
func handleDel(args []string, conn net.Conn, storage *server.MStorageServer) {
	if len(args) != 2 {
		conn.Write([]byte("-ERR wrong number of arguments for 'DEL' command\r\n"))
		return
	}
	key := args[1]
	if storage.Del(key) {
		conn.Write([]byte(":1\r\n"))
	} else {
		conn.Write([]byte(":0\r\n"))
	}
}
