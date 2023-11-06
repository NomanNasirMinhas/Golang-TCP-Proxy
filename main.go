package main

import (
	"flag"
)

func main() {
	// Get the address and port to bind to, and whether to start a server or client
	var addr string
	var isServer bool
	var isClient bool
	flag.StringVar(&addr, "addr", "localhost:8080", "Address of the server")
	flag.BoolVar(&isServer, "server", false, "Start in server mode")
	flag.BoolVar(&isClient, "client", false, "Start in client mode")
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
		startServer(addr)
	} else {
		startClient(addr)
	}
}
