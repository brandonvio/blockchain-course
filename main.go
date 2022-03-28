package main

import (
	"context"
	"go.uber.org/fx"
	"log"
	"time"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	app := fx.New(
		fx.Provide(NewGlobals),
		fx.Provide(NewBlockchain),
		fx.Invoke(RunBlockChain),
	)
	startCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := app.Start(startCtx); err != nil {
		log.Fatal(err)
	}
}

func RunBlockChain(bc *Blockchain) {
	bc.AddTransaction("A", "B", 2.2)
	bc.CreateBlock(5)

	bc.AddTransaction("C", "D", 3.1)
	bc.CreateBlock(2)

	bc.AddTransaction("X", "Y", 3.1)
	bc.AddTransaction("J", "K", 4.5)
	bc.CreateBlock(2)
	bc.Print()
}
