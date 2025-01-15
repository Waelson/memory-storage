package resp

import (
	"net"

	"github.com/Waelson/memory-storage/m-storage/internal/server"
)

// handleFlushAll processa o comando FLUSHALL.
func handleFlushAll(args []string, conn net.Conn, storage *server.MStorageServer) {
	if len(args) != 1 {
		conn.Write([]byte("-ERR wrong number of arguments for 'FLUSHALL' command\r\n"))
		return
	}

	storage.LockData()
	defer storage.UnlockData()

	for key := range storage.GetData() {
		delete(storage.GetData(), key)
	}

	// Persistir a limpeza em disco
	storage.SignalWrite(storage.GetData())

	conn.Write([]byte("+OK\r\n"))
}
