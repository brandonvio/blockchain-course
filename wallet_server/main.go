package main

import (
	"flag"
	"log"
)

func init() {
	log.SetPrefix("wallet_server:")
}

func main() {
	port := flag.Uint("port", 8080, "TCP Port for Wallet Server")
	gateway := flag.String("gateway", "http://127.0.0.1:5000", "Blockchain Gateway")
	flag.Parse()
	ws := NewWalleteServer(*port, *gateway)
	ws.Run()
}
