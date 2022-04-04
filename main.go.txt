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
		fx.Invoke(CreateWallets),
	)
	startCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := app.Start(startCtx); err != nil {
		log.Fatal(err)
	}
}

func CreateWallets(bc *block.Blockchain) {
	walletMiner := wallet.NewWallet()
	walletA := wallet.NewWallet()
	walletB := wallet.NewWallet()

	bc.SetBlockchainAddress(walletMiner.BlockchainAddress())

	t1 := wallet.NewTransaction(
		walletA.PrivateKey(),
		walletA.PublicKey(),
		walletA.BlockchainAddress(),
		walletB.BlockchainAddress(),
		2.0,
	)

	t1signature := t1.GenerateSignature()

	isAdded := bc.AddTransaction(
		walletA.BlockchainAddress(),
		walletB.BlockchainAddress(),
		2.0,
		walletA.PublicKey(),
		t1signature,
	)

	fmt.Println("Added? ", isAdded)
	bc.Mining()
	bc.Print()
}

func RunBlockChain(bc *block.Blockchain) {
	// bc.AddTransaction("A", "B", 2.2)
	// bc.Mining()

	// bc.AddTransaction("C", "B", 3.1)
	// bc.Mining()

	// bc.AddTransaction("X", "Y", 3.1)
	// bc.AddTransaction("J", "K", 4.5)
	// bc.AddTransaction("B", "K", 1.2)
	// bc.Mining()
	// bc.Print()

	// fmt.Printf("miner	%.2f\n", bc.CalculateTotalAmount("my_blockchain_address"))
	// fmt.Printf("A	%.2f\n", bc.CalculateTotalAmount("A"))
	// fmt.Printf("B	%.2f\n", bc.CalculateTotalAmount("B"))
}
