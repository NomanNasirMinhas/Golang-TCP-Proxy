package main

import (
	"flag"
)

func main() {
	// Get the address and port to bind to, and whether to start a server or client
	var addr string
	var isServer bool
	var isClient bool
	var useTLS bool
	var certFile string
	var keyFile string
	var generateCert bool
	flag.StringVar(&addr, "addr", "0.0.0.0:8080", "Address of the server")
	flag.BoolVar(&isServer, "server", false, "Start in server mode")
	flag.BoolVar(&isClient, "client", false, "Start in client mode")
	flag.BoolVar(&useTLS, "tls", false, "Use TLS")
	flag.StringVar(&certFile, "cert", "cert.pem", "Certificate file")
	flag.StringVar(&keyFile, "key", "key.pem", "Key file")
	flag.BoolVar(&generateCert, "generate-cert", false, "Generate a self-signed certificate")
	flag.Parse()
	// show help if no arguments are passed
	if flag.NFlag() == 0 {
		flag.Usage()
		return
	}
	if isServer && isClient {
		panic("Cannot start in both server and client mode")
	}

	if isServer {
		startServer(addr, useTLS, certFile, keyFile, generateCert)
	} else {
		startClient(addr, useTLS, certFile)
	}
}
