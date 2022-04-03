package main

import (
	types "blockchain/blockchaintypes"

	"fmt"
	"strings"
)

const DIFFICULTY_LEVEL = 3

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

func (bc *Blockchain) CreateBlock() *Block {
	lb := bc.LastBlock()
	previousHash := lb.Hash()
	nonce := bc.ProofOfWork()
	b := NewBlock(nonce, previousHash, bc.globals.NowUnixNano(), bc.transactionPool)
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

func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, transaction := range bc.transactionPool {
		transactions = append(transactions,
			NewTransaction(
				transaction.sendBlockchainAddress,
				transaction.recipientBlockchainAddress,
				transaction.value))
	}
	return transactions
}

func (bc *Blockchain) ValidProof(
	nonce int,
	previousHash types.Byte32,
	transactions []*Transaction,
	difficulty int) bool {

	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{
		timestamp:    0,
		nonce:        nonce,
		previousHash: previousHash,
		transactions: transactions,
	}
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	// fmt.Println(guessHashStr)
	return guessHashStr[:difficulty] == zeros
}

func (bc *Blockchain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	previousHash := bc.LastBlock().Hash()
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transactions, DIFFICULTY_LEVEL) {
		nonce += 1
	}
	return nonce
}
