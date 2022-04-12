package main

import (
	"blockchain/block"
	"blockchain/globals"
	"blockchain/wallet"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
)

var gl = globals.NewGlobals()

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

func (bcs *BlockchainServer) Transactions(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := bcs.GetBlockchain()
		transactions := bc.TransactionPool()
		m, _ := json.Marshal(struct {
			Transactions []*block.Transaction `json:"transactions"`
			Length       int                  `json:"length"`
		}{
			Transactions: transactions,
			Length:       len(transactions),
		})
		io.WriteString(w, string(m[:]))

	// TODO
	case http.MethodPost:
		decoder := json.NewDecoder(req.Body)
		var t block.TransactionRequest
		err := decoder.Decode(&t)
		if err != nil {
			log.Printf("ERROR: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, string(gl.JsonStatus("failed: decoding failed")))
			return
		}
		if !t.Validate() {
			log.Printf("ERROR: invalid payload")
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, string(gl.JsonStatus("failed: invalid payload")))
		}

		publicKey := gl.PublicKeyFromString(*t.SenderPublicKey)
		signature := gl.SignatureFromString(*t.Signature)

		isCreated := bcs.blockchain.CreateTransaction(
			*t.SenderBlockchainAddress,
			*t.RecipientBlockchainAddress,
			*t.Value,
			publicKey,
			signature,
		)

		w.Header().Add("Content-Type", "application/json")
		var m []byte
		if !isCreated {
			w.WriteHeader(http.StatusBadRequest)
			m = gl.JsonStatus("fail")
		} else {
			w.WriteHeader(http.StatusCreated)
			m = gl.JsonStatus("success")
		}
		io.WriteString(w, string(gl.JsonStatus(string(m))))

		log.Printf("INFO: transaction_request: %+v", t)
	default:
		log.Printf("ERROR: invalid http method")
	}
}

func (bcs *BlockchainServer) Run(port uint16) {
	bcs.port = port
	bcs.blockchain.SetPort(port)
	log.Printf("Starting blockchain server with port %v", port)
	http.HandleFunc("/", bcs.GetChain)
	http.HandleFunc("/transactions", bcs.Transactions)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(bcs.Port())), nil))
}

// func (bcs *BlockchainServer) GetBlockchain() *block.Blockchain {
// 	bc, ok := cache["blockchain"]
// 	if !ok {
// 		minersWallet := wallet.NewWallet()
// 		bc = block.NewBlockchain()
// 	}
// }
