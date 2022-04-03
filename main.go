package main

import (
	"blockchain/block"
	"blockchain/globals"
	"blockchain/wallet"
	"context"
	"fmt"
	"log"
	"time"

	"go.uber.org/fx"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	app := fx.New(
		fx.Provide(globals.NewGlobals),
		fx.Provide(block.NewBlockchain),
		fx.Invoke(CreateWallet),
	)
	startCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := app.Start(startCtx); err != nil {
		log.Fatal(err)
	}
}

func CreateWallet() {
	w := wallet.NewWallet()
	fmt.Println(w.PrivateKeyStr())
	fmt.Println(w.PublicKeyStr())
	fmt.Println(w.BlockchainAddress())
}

func RunBlockChain(bc *block.Blockchain) {
	bc.AddTransaction("A", "B", 2.2)
	bc.Mining()

	bc.AddTransaction("C", "B", 3.1)
	bc.Mining()

	bc.AddTransaction("X", "Y", 3.1)
	bc.AddTransaction("J", "K", 4.5)
	bc.AddTransaction("B", "K", 1.2)
	bc.Mining()
	bc.Print()

	fmt.Printf("miner	%.2f\n", bc.CalculateTotalAmount("my_blockchain_address"))
	fmt.Printf("A	%.2f\n", bc.CalculateTotalAmount("A"))
	fmt.Printf("B	%.2f\n", bc.CalculateTotalAmount("B"))
}
