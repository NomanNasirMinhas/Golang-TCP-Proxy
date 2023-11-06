package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os"
)

func startClient(addr string, useTLS bool, certFile string) {
	var conn net.Conn
	var err error
	// Connect to the server
	if useTLS {
		cert, err := os.ReadFile(certFile)
		if err != nil {
			log.Fatal(err)
		}
		certPool := x509.NewCertPool()
		if ok := certPool.AppendCertsFromPEM(cert); !ok {
			log.Fatalf("unable to parse cert from %s", certFile)
		}
		config := &tls.Config{RootCAs: certPool}

		conn, err = tls.Dial("tcp", addr, config)
	} else {
		conn, err = net.Dial("tcp", addr)
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	// Send some data to the server
	_, err = conn.Write([]byte("Hello, server!"))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Close the connection
	conn.Close()
}
