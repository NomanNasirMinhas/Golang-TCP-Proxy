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
	// Start an infinite goroutine to read from the server and print to stdout
	go func() {
		for {
			// Read from the server
			buf := make([]byte, 1024*10)
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println(err)
				return
			}
			// Print to stdout
			fmt.Print("Server Sent: " + string(buf[:n]) + "\n")
		}
	}()

	// Go into an infinite loop reading from stdin and sending to the server
	for {
		// Read from stdin
		buf := make([]byte, 1024*10)
		n, err := os.Stdin.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}

		//If the user types "quit", then stop
		if string(buf[:n-1]) == "quit" {
			break
		}
		// Send to the server
		_, err = conn.Write(buf[:n])
		if err != nil {
			fmt.Println(err)
			return
		}

	}

	// Close the connection
	conn.Close()
}
