package main

import (
	"context"
	"log"
	"time"

	"go.uber.org/fx"
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
	bc.CreateBlock()

	bc.AddTransaction("C", "D", 3.1)
	bc.CreateBlock()

	bc.AddTransaction("X", "Y", 3.1)
	bc.AddTransaction("J", "K", 4.5)
	bc.CreateBlock()
	bc.Print()
}
