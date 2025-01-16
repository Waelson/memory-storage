package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"github.com/Waelson/memory-storage/m-storage/internal/resp"
	"github.com/Waelson/memory-storage/m-storage/internal/server"
	"io/ioutil"
	"log"
	"net"
	"os"
)

func main() {
	displayLogo()

	// Flags para configurar o servidor
	var (
		enableTLS    = flag.Bool("tls", false, "Enable TLS for secure connections (default: false)")
		certPath     = flag.String("cert", "certs/server.crt", "Path to the server certificate")
		keyPath      = flag.String("key", "certs/server.key", "Path to the server private key")
		caCertPath   = flag.String("ca", "certs/ca.crt", "Path to the CA certificate")
		insecurePort = flag.String("insecure-port", "6379", "Port for insecure connections (default: 6379)")
		securePort   = flag.String("secure-port", "6380", "Port for secure connections (only used if --tls is set)")
	)
	flag.Parse()

	filePath := "mstorage.dat"
	svr := server.NewRedisServer(filePath)

	// Sempre inicia o servidor não seguro
	go startInsecureServer(svr, *insecurePort)

	// Inicia o servidor seguro apenas se TLS estiver habilitado
	if *enableTLS {
		validateTLSFiles(*certPath, *keyPath, *caCertPath)
		go startSecureServer(svr, *securePort, *certPath, *keyPath, *caCertPath)
	}

	// Mantém o programa rodando
	select {}
}

// startInsecureServer inicia o servidor em uma porta não segura
func startInsecureServer(svr *server.MStorageServer, port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to start insecure server: %v", err)
	}
	fmt.Printf("M-Storage insecure server is running on :%s\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Connection error: %v", err)
			continue
		}
		go resp.HandleRESP(conn, svr)
	}
}

// startSecureServer inicia o servidor em uma porta segura com TLS
func startSecureServer(svr *server.MStorageServer, port, certPath, keyPath, caCertPath string) {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Fatalf("Failed to load TLS certificates: %v", err)
	}

	caCert, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		log.Fatalf("Failed to read CA certificate: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	listener, err := tls.Listen("tcp", ":"+port, tlsConfig)
	if err != nil {
		log.Fatalf("Failed to start secure server: %v", err)
	}
	fmt.Printf("M-Storage secure server is running on :%s\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("TLS connection error: %v", err)
			continue
		}
		go resp.HandleRESP(conn, svr)
	}
}

// validateTLSFiles verifica se os arquivos de certificado e chave existem
func validateTLSFiles(certPath, keyPath, caCertPath string) {
	files := []string{certPath, keyPath, caCertPath}
	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			log.Fatalf("Required TLS file not found: %s", file)
		}
	}
}

func displayLogo() {
	logo := `
        ,MMM8&&&.            *
   _...MMMMM88&&&&...       .
 .::'''MMMMM88&&&&&&'''::. 
::     MMMMM88&&&&&&     ::  
'::....MMMMM88&&&&&&....::'  
   ` + "`" + `::::MMMMM88&&&&::::'      
       ` + "`" + `::::88888::::'          
           ` + "`" + `:::::'            

      Saturn - M-Storage
	`
	fmt.Println(logo)
}
