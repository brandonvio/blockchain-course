package main

import (
	"blockchain/globals"
	"blockchain/wallet"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"
	"text/template"
)

const tempDir = "templates"

type WalletServer struct {
	port    uint
	gateway string
	lib     globals.IGlobalLib
}

type Transaction struct {
	SenderPrivateKey                 string `json:"sender_private_key"`
	SenderBlockchainAddress          string `json:"sender_blockchain_address"`
	SenderRecipientBlockchainAddress string `json:"sender_recipient_blockchain_address"`
	SenderPublicKey                  string `json:"sender_public_key"`
	SenderSendAmount                 string `json:"sender_send_amount"`
}

func NewWalletServer(port uint, gateway string, lib globals.IGlobalLib) *WalletServer {
	return &WalletServer{
		port:    port,
		gateway: gateway,
		lib:     lib,
	}
}

func (ws *WalletServer) Port() uint {
	return ws.port
}

func (ws *WalletServer) Gateway() string {
	return ws.gateway
}

func (ws *WalletServer) Index(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		t, _ := template.ParseFiles(path.Join(tempDir, "index.html"))
		err := t.Execute(w, "")
		if err != nil {
			log.Fatalln(err)
		}
	default:
		log.Printf("ERROR: Invalid HTTP Method")
	}
}

func (ws *WalletServer) Wallet(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		w.Header().Add("Content-Type", "application/json")
		myWallet := wallet.NewWallet()
		m, _ := myWallet.MarshalJSON()
		_, err := w.Write(m)
		if err != nil {
			log.Println(err)
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method")
	}
}

func (ws *WalletServer) CreateTransaction(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		// handle request
		var tx Transaction
		err := ws.lib.DecodeJSONBody(w, req, &tx)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("tx.SenderPrivateKey:", tx.SenderPrivateKey)

		// send response
		w.Header().Add("Content-Type", "application/json")
		_, err = io.WriteString(w, string(ws.lib.JsonStatus("success")))
		if err != nil {
			log.Println(err)
		}
		log.Println("create transaction success")
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method")
	}
}

func (ws *WalletServer) Run() {
	http.HandleFunc("/wallet", ws.Wallet)
	http.HandleFunc("/transaction", ws.CreateTransaction)
	http.HandleFunc("/", ws.Index)
	log.Printf("Running wallet server on port %v\n", ws.Port())
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(ws.Port())), nil))
}
