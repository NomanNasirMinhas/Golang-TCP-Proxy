package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/open-quantum-safe/liboqs-go/oqs"
)

func startServer(addr string, TLS_Type int8, certFile string, keyFile string) {
	d := oqs.LiboqsVersion()
	fmt.Printf("Lib-OQS version: %s\n", d)

	var ln net.Listener
	var err error
	if TLS_Type == 1 {
		println("Running TLS server")
		// Check if the certificate and key files exist
		if _, err := os.Stat(certFile); os.IsNotExist(err) {
			log.Printf("certificate file %s does not exist. Generating", certFile)
			generateCertFunc(strings.Split(addr, ":")[0])
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
			log.Fatalf("certificate file %s does not exist.", certFile)
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
	api_url := "https://randomuser.me/api/?seed="
	// Go into an infinite loop reading from the connection
	for {
		// Read from the connection
		buf := make([]byte, 1024*10)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		// If the user types "quit", then stop
		if string(buf[:n-1]) == "quit" {
			break
		}
		query := api_url + string(buf[:n-1])
		fmt.Printf("Performing a random query for %s\n", query)
		//Define a transport to object to set the oqs ciphersuite
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					//CipherSuites: []uint16{oqs.TLS_AES_128_OQS_KEM_FRODOKEM_640_AES_128_SHA256},
				},
			}}
		// get request on the api
		resp, err := client.Get(strings.Split(query, "\r")[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		// read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		// write the response body to the connection
		_, err = conn.Write(body)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Sent to %s: \n", conn.RemoteAddr())

	}

	// Print the incoming data
}
