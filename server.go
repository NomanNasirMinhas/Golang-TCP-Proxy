package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func startServer(addr string, TLS_Type int8, certFile string, keyFile string, generateCert bool) {
	var ln net.Listener
	var err error
	if TLS_Type == 1 {
		println("Running TLS server")
		// Check if the certificate and key files exist
		if _, err := os.Stat(certFile); os.IsNotExist(err) {
			log.Printf("certificate file %s does not exist.", certFile)
			if generateCert {
				log.Printf("generating self-signed certificate %s", certFile)
				generateCertFunc(strings.Split(addr, ":")[0])
			} else {
				log.Fatalf("use -generate-cert to generate a self-signed certificate")
			}
		}
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Fatal(err)
		}
		config := &tls.Config{Certificates: []tls.Certificate{cert}}
		log.Printf("listening on port %s\n", addr)
		ln, err = tls.Listen("tcp", addr, config)
		if err != nil {
			fmt.Println(err)
			return
		}

	} else if TLS_Type == 2 {
		println("Running PQ-TLS server")
		// Check if the certificate and key files exist
		if _, err := os.Stat(certFile); os.IsNotExist(err) {
			log.Printf("certificate file %s does not exist.", certFile)
			if generateCert {
				log.Printf("generating self-signed certificate %s", certFile)
				generateCertFunc(strings.Split(addr, ":")[0])
			} else {
				log.Fatalf("use -generate-cert to generate a self-signed certificate")
			}
		}
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Fatal(err)
		}
		config := &tls.Config{Certificates: []tls.Certificate{cert}}
		log.Printf("listening on port %s\n", addr)
		ln, err = tls.Listen("tcp", addr, config)
		if err != nil {
			fmt.Println(err)
			return
		}

	} else {
		println("Running plaintext server")
		ln, err = net.Listen("tcp", addr)
		if err != nil {
			fmt.Println(err)
			return
		}
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

func handleConnection(conn net.Conn) {
	// Close the connection when we're done
	defer conn.Close()

	// Go into an infinite loop reading from the connection
	for {
		// Read from the connection
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		// If the user types "quit", then stop
		if string(buf[:n-1]) == "quit" {
			break
		}

		fmt.Printf("Received from %s: %s\n", conn.RemoteAddr(), buf)
	}

	// Print the incoming data
}
