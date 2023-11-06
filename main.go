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
	var usePQC bool
	var certFile string
	var keyFile string
	var generateCert bool
	flag.StringVar(&addr, "addr", "0.0.0.0:8080", "Address of the server")
	flag.BoolVar(&isServer, "server", false, "Start in server mode")
	flag.BoolVar(&isClient, "client", false, "Start in client mode")
	flag.BoolVar(&useTLS, "tls", false, "Use TLS")
	flag.BoolVar(&usePQC, "pq_tls", false, "Use PQ-TLS")
	flag.StringVar(&certFile, "cert", "cert.pem", "Certificate file")
	flag.StringVar(&keyFile, "key", "key.pem", "Key file")
	flag.BoolVar(&generateCert, "generate-cert", false, "Generate a self-signed certificate")
	flag.Parse()
	// show help if no arguments are passed
	if flag.NFlag() == 0 {
		flag.Usage()
		return
	}
	if !isServer && !isClient {
		panic("Must start in either server or client mode")
	}
	if isServer && isClient {
		panic("Cannot start in both server and client mode")
	}
	if usePQC && useTLS {
		panic("Cannot use both TLS and PQ-TLS")
	}
	tls_type := 0
	if useTLS {
		tls_type = 1
	} else if usePQC {
		tls_type = 2
	} else {
		tls_type = 0
	}

	if isServer {
		startServer(addr, int8(tls_type), certFile, keyFile, generateCert)
	} else {
		startClient(addr, useTLS, certFile)
	}
}
