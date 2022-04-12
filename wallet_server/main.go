package main

import (
	"blockchain/globals"
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

	log.Println("INFO: Blockchain gateway configured as:", *gateway)
	lib := &globals.GlobalLib{}

	ws := NewWalletServer(*port, *gateway, lib)
	ws.Run()
}
