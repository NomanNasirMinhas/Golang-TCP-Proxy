package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func startServer(addr string, useTLS bool, certFile string, keyFile string, generateCert bool) {
	var ln net.Listener
	var err error
	if useTLS {
		// Check if the certificate and key files exist
		if _, err := os.Stat(certFile); os.IsNotExist(err) {
			log.Printf("certificate file %s does not exist.", certFile)
			if generateCert {
				log.Printf("generating self-signed certificate %s", certFile)
				generateCertFunc(strings.Split(addr, ":")[0])
			} else {
				log.Fatalf("use -generate-cert to generate a self-signed certificate")
			}
			if _, err := os.Stat(keyFile); os.IsNotExist(err) {
				log.Fatalf("key file %s does not exist", keyFile)
			}
			cert, err := tls.LoadX509KeyPair(certFile, keyFile)
			if err != nil {
				log.Fatal(err)
			}
			config := &tls.Config{Certificates: []tls.Certificate{cert}}
			log.Printf("listening on port %s\n", addr)
			ln, err = tls.Listen("tcp", addr, config)
		} else {
			ln, err = net.Listen("tcp", addr)
		}
		if err != nil {
			fmt.Println(err)
			return
		}

		// Accept incoming connections and handle them
		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println(err)
				continue
			}

			// Handle the connection in a new goroutine
			go handleConnection(conn)
		}
	}
}

func handleConnection(conn net.Conn) {
	// Close the connection when we're done
	defer conn.Close()

	// Read incoming data
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the incoming data
	fmt.Printf("Received: %s", buf)
}
