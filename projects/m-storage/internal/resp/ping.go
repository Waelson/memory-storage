package resp

import (
	"fmt"
	"net"
)

// handlePing processa o comando PING.
func handlePing(args []string, conn net.Conn) {
	if len(args) == 1 {
		conn.Write([]byte("+PONG\r\n")) // Resposta padrão
	} else if len(args) == 2 {
		conn.Write([]byte(fmt.Sprintf("+%s\r\n", args[1]))) // Retorna o argumento
	} else {
		conn.Write([]byte("-ERR wrong number of arguments for 'PING' command\r\n"))
	}
}
