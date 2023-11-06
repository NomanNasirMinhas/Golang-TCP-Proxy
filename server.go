package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
)

func startServer(addr string, useTLS bool, certFile string, keyFile string) {
	var ln net.Listener
	var err error
	if useTLS {
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
