package resp

import (
	"bufio"
	"fmt"
	"github.com/Waelson/memory-storage/m-storage/internal/server"
	"io"
	"net"
	"strconv"
	"strings"
)

// HandleRESP lida com requisições RESP e delega os comandos.
func HandleRESP(conn net.Conn, storage *server.MStorageServer) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client disconnected.")
				return
			}
			fmt.Println("Error reading request:", err)
			return
		}

		// Processa o número de argumentos do protocolo RESP
		if len(line) == 0 || line[0] != '*' {
			conn.Write([]byte("-ERR invalid command\r\n"))
			continue
		}

		numArgs, err := strconv.Atoi(strings.TrimSpace(line[1:]))
		if err != nil {
			conn.Write([]byte("-ERR invalid number of arguments\r\n"))
			continue
		}

		// Leitura dos argumentos
		args := make([]string, 0, numArgs)
		for i := 0; i < numArgs; i++ {
			_, err = reader.ReadString('\n') // Lê o tamanho do argumento (não usado)
			if err != nil {
				conn.Write([]byte("-ERR invalid argument\r\n"))
				return
			}
			arg, err := reader.ReadString('\n')
			if err != nil {
				conn.Write([]byte("-ERR invalid argument\r\n"))
				return
			}
			args = append(args, strings.TrimSpace(arg))
		}

		// Processa o comando
		if len(args) < 1 {
			conn.Write([]byte("-ERR no command provided\r\n"))
			continue
		}

		cmd := strings.ToUpper(args[0])
		executeCommand(cmd, args, conn, storage)
	}
}

// executeCommand delega o comando para a função apropriada.
func executeCommand(cmd string, args []string, conn net.Conn, storage *server.MStorageServer) {
	switch cmd {
	case "PING":
		handlePing(args, conn)
	case "SET":
		handleSet(args, conn, storage)
	case "GET":
		handleGet(args, conn, storage)
	case "DEL":
		handleDel(args, conn, storage)
	case "EXISTS":
		handleExists(args, conn, storage)
	case "EXPIRE":
		handleExpire(args, conn, storage)
	case "TTL":
		handleTTL(args, conn, storage)
	case "FLUSHALL":
		handleFlushAll(args, conn, storage)
	case "INFO":
		handleInfo(args, conn, storage)
	case "QUIT":
		conn.Write([]byte("+OK\r\n"))
		fmt.Println("Client requested to close the connection.")
		conn.Close()
	default:
		conn.Write([]byte("-ERR unknown command\r\n"))
	}
}
