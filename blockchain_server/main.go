package main

import (
	"blockchain/block"
	"blockchain/globals"
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"go.uber.org/fx"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func StartServer(bcs *BlockchainServer) {
	port := flag.Uint("port", 5000, "TCP Port Number for Blockchain Server")
	flag.Parse()
	fmt.Println(*port)
	bcs.Run(uint16(*port))
}

func main() {
	app := fx.New(
		fx.Provide(globals.NewGlobals),
		fx.Provide(block.NewBlockchain),
		fx.Provide(NewBlockchainServer),
		fx.Invoke(StartServer),
	)
	startCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := app.Start(startCtx); err != nil {
		log.Fatal(err)
	}
}
