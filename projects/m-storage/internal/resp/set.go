package resp

import (
	"net"
	"strconv"
	"time"

	"github.com/Waelson/memory-storage/m-storage/internal/server"
)

// handleSet processa o comando SET.
func handleSet(args []string, conn net.Conn, storage *server.MStorageServer) {
	if len(args) != 4 {
		conn.Write([]byte("-ERR wrong number of arguments for 'SET' command\r\n"))
		return
	}
	key, value := args[1], args[2]
	ttl, err := strconv.Atoi(args[3])
	if err != nil {
		conn.Write([]byte("-ERR invalid TTL\r\n"))
		return
	}
	storage.Set(key, value, time.Duration(ttl)*time.Second)
	conn.Write([]byte("+OK\r\n"))
}
