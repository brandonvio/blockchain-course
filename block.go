package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"

	types "blockchain/blockchaintypes"
)

type Block struct {
	nonce        int
	previousHash types.Byte32
	timestamp    int64
	transactions []*Transaction
}

func NewBlock(nonce int, previousHash types.Byte32, timestamp int64, transactions []*Transaction) *Block {
	return &Block{
		nonce:        nonce,
		previousHash: previousHash,
		timestamp:    timestamp,
		transactions: transactions,
	}
}

func (b *Block) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf("timestamp        %d\n", b.timestamp)
	fmt.Printf("nonce            %d\n", b.nonce)
	fmt.Printf("previousHash     %x\n", b.previousHash)
	fmt.Printf("hash             %x\n", b.Hash())
	for _, t := range b.transactions {
		t.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("-", 40))
}

func (b *Block) Hash() types.Byte32 {
	m, _ := json.Marshal(b)
	return sha256.Sum256(m)
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64          `json:"timestamp"`
		Nonce        int            `json:"nonce"`
		PreviousHash types.Byte32   `json:"previous_hash"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Timestamp:    b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: b.previousHash,
		Transactions: b.transactions,
	})
}
