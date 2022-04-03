package block

import (
	types "blockchain/blockchaintypes"
	"blockchain/globals"
	"log"

	"fmt"
	"strings"
)

const (
	MINING_DIFFICULTY = 3
	MINING_SENDER     = "THE BLOCKCHAIN"
	MINING_REWARD     = 1.0
)

type Blockchain struct {
	globals           globals.IGlobalLib
	transactionPool   []*Transaction
	chain             []*Block
	blockchainAddress string
}

func NewBlockchain(globals globals.IGlobalLib) *Blockchain {
	bc := new(Blockchain)
	bc.blockchainAddress = "my_blockchain_address"
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
	return guessHashStr[:difficulty] == zeros
}

func (bc *Blockchain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	previousHash := bc.LastBlock().Hash()
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transactions, MINING_DIFFICULTY) {
		nonce += 1
	}
	return nonce
}

func (bc *Blockchain) Mining() bool {
	bc.AddTransaction(MINING_SENDER, bc.blockchainAddress, MINING_REWARD)
	bc.CreateBlock()
	log.Println("action=mining status=success")
	return true
}

func (bc *Blockchain) CalculateTotalAmount(blockchainAddress string) float32 {
	var totalAmount float32 = 0.0
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			value := t.value
			if t.recipientBlockchainAddress == blockchainAddress {
				totalAmount += value
			}

			if t.sendBlockchainAddress == blockchainAddress {
				totalAmount -= value
			}
		}
	}
	return totalAmount
}
