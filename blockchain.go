package main

import (
	"fmt"
	"strings"
)

type Blockchain struct {
	globals         IGlobalLib
	transactionPool []*Transaction
	chain           []*Block
}

func NewBlockchain(globals IGlobalLib) *Blockchain {
	bc := new(Blockchain)
	bc.globals = globals
	b0 := NewBlock(0, globals.EmptyByte32(), globals.NowUnixNano(), []*Transaction{})
	bc.chain = append(bc.chain, b0)
	return bc
}

func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Block %d %s\n", strings.Repeat("=", 15), i, strings.Repeat("=", 15))
		block.Print()
		fmt.Printf("%s\n", strings.Repeat("~", 39))
	}
	fmt.Printf("%s\n", strings.Repeat("*", 39))
}

func (bc *Blockchain) CreateBlock(nonce int) *Block {
	lb := bc.LastBlock()
	b := NewBlock(nonce, lb.Hash(), bc.globals.NowUnixNano(), bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	return b
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) AddTransaction(sender string, recipient string, value float32) *Transaction {
	t := NewTransaction(sender, recipient, value)
	bc.transactionPool = append(bc.transactionPool, t)
	return t
}
