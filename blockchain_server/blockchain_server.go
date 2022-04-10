package main

import (
	"blockchain/block"
	"blockchain/wallet"
	"io"
	"log"
	"net/http"
	"strconv"
)

type BlockchainServer struct {
	port       uint16
	blockchain *block.Blockchain
}

func NewBlockchainServer(blockchain *block.Blockchain) *BlockchainServer {
	minersWallet := wallet.NewWallet()
	blockchain.SetBlockchainAddress(minersWallet.BlockchainAddress())
	log.Printf("miner's private_key %v", minersWallet.PrivateKeyStr())
	log.Printf("miner's public_key %v", minersWallet.PublicKeyStr())
	log.Printf("miner's blockchain_address %v", minersWallet.BlockchainAddress())
	return &BlockchainServer{
		blockchain: blockchain,
	}
}

func (bcs *BlockchainServer) GetBlockchain() *block.Blockchain {
	return bcs.blockchain
}

func (bcs *BlockchainServer) Port() uint16 {
	return bcs.port
}

func (bcs *BlockchainServer) GetChain(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := bcs.GetBlockchain()
		m, _ := bc.MarshalJSON()
		_, err := io.WriteString(w, string(m[:]))
		if err != nil {
			log.Fatalln(err)
			return
		}
	default:
		log.Printf("ERROR: Invalid HTTP Method")
	}
}

func (bcs *BlockchainServer) Run(port uint16) {
	bcs.port = port
	bcs.blockchain.SetPort(port)
	log.Printf("Starting blockchain server with port %v", port)
	http.HandleFunc("/", bcs.GetChain)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(bcs.Port())), nil))
}

// func (bcs *BlockchainServer) GetBlockchain() *block.Blockchain {
// 	bc, ok := cache["blockchain"]
// 	if !ok {
// 		minersWallet := wallet.NewWallet()
// 		bc = block.NewBlockchain()
// 	}
// }
