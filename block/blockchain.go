package block

import (
	types "blockchain/blockchaintypes"
	"blockchain/globals"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"log"
	"sync"
	"time"

	"fmt"
	"strings"
)

const (
	MiningDifficulty = 3
	MiningSender     = "THE BLOCKCHAIN"
	MiningReward     = 1.0
	MiningTimerSec   = 20

	BlockchainPortRangeStart      = 5001
	BlockchainPortRangeEnd        = 5004
	NeighborIpRangeStart          = 0
	NeighborIpRangeEnd            = 1
	BlockchainNeighborSyncTimeSec = 20
)

type Blockchain struct {
	globals           globals.IGlobalLib
	transactionPool   []*Transaction
	chain             []*Block
	blockchainAddress string
	port              uint16
	mux               sync.Mutex
	neighbors         []string
	muxNeighbors      sync.Mutex
}

type AmountResponse struct {
	Amount float32 `json:"amount"`
}

func (ar *AmountResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Amount float32 `json:"amount"`
	}{
		Amount: ar.Amount,
	})
}

func NewBlockchain(globals globals.IGlobalLib) *Blockchain {
	bc := new(Blockchain)
	bc.blockchainAddress = "my_blockchain_address"
	bc.globals = globals
	b0 := NewBlock(0, globals.EmptyByte32(), globals.NowUnixNano(), []*Transaction{})
	bc.chain = append(bc.chain, b0)
	return bc
}

func (bc *Blockchain) Run() {
	bc.StartSyncNeighbors()
}

func (bc *Blockchain) SetNeighbors() {
	myHost := globals.GetHost()
	bc.neighbors = globals.FindNeighbors(
		myHost,
		bc.port,
		NeighborIpRangeStart,
		NeighborIpRangeEnd,
		BlockchainPortRangeStart,
		BlockchainPortRangeEnd,
	)
}

func (bc *Blockchain) SyncNeighbors() {
	bc.muxNeighbors.Lock()
	defer bc.muxNeighbors.Unlock()
	bc.SetNeighbors()
}

func (bc *Blockchain) StartSyncNeighbors() {
	bc.SyncNeighbors()
	_ = time.AfterFunc(time.Second*BlockchainNeighborSyncTimeSec, bc.StartSyncNeighbors)
}

func (bc *Blockchain) TransactionPool() []*Transaction {
	return bc.transactionPool
}

func (bc *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*Block `json:"blocks"`
	}{
		Blocks: bc.chain,
	})
}

func (bc *Blockchain) SetBlockchainAddress(address string) {
	bc.blockchainAddress = address
}

func (bc *Blockchain) SetPort(port uint16) {
	bc.port = port
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
	timestamp := bc.globals.NowUnixNano()
	b := NewBlock(nonce, previousHash, timestamp, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	return b
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) VerifyTransactionSignature(
	senderPublicKey *ecdsa.PublicKey,
	s *globals.Signature,
	t *Transaction,
) bool {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	return ecdsa.Verify(senderPublicKey, h[:], s.R, s.S)
}

func (bc *Blockchain) CreateTransaction(
	sender string,
	recipient string,
	value float32,
	senderPublicKey *ecdsa.PublicKey,
	s *globals.Signature) bool {

	isTransacted := bc.AddTransaction(
		sender,
		recipient,
		value,
		senderPublicKey,
		s)

	return isTransacted
}

func (bc *Blockchain) AddTransaction(
	sender string,
	recipient string,
	value float32,
	senderPublicKey *ecdsa.PublicKey,
	s *globals.Signature) bool {

	t := NewTransaction(sender, recipient, value)

	if sender == MiningSender {
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}

	if bc.VerifyTransactionSignature(senderPublicKey, s, t) {

		// if bc.CalculateTotalAmount(sender) < value {
		// 	log.Println("ERROR: not enough balance in wallet")
		// 	return false
		// }

		bc.transactionPool = append(bc.transactionPool, t)
		return true
	} else {
		log.Println("ERROR: Verify Transaction")
		return false
	}
}

func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, transaction := range bc.transactionPool {
		transactions = append(transactions,
			NewTransaction(
				transaction.senderBlockchainAddress,
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
	for !bc.ValidProof(nonce, previousHash, transactions, MiningDifficulty) {
		nonce += 1
	}
	return nonce
}

func (bc *Blockchain) Mining() bool {
	bc.mux.Lock()
	defer bc.mux.Unlock()

	if len(bc.transactionPool) > 0 {
		bc.AddTransaction(MiningSender, bc.blockchainAddress, MiningReward, nil, nil)
		bc.CreateBlock()
		log.Println("action=mining status=success")
		return true
	} else {
		log.Println("action=mining status=zero transactions to mine")
		return false
	}
}

func (bc *Blockchain) StartMining() {
	bc.Mining()
	_ = time.AfterFunc(time.Second*MiningTimerSec, bc.StartMining)
}

func (bc *Blockchain) CalculateTotalAmount(blockchainAddress string) float32 {
	var totalAmount float32 = 0.0
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			value := t.value
			if t.recipientBlockchainAddress == blockchainAddress {
				totalAmount += value
			}

			if t.senderBlockchainAddress == blockchainAddress {
				totalAmount -= value
			}
		}
	}
	return totalAmount
}
