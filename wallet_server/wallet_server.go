package main

import (
	"blockchain/block"
	"blockchain/globals"
	"blockchain/wallet"
	"bytes"
	"encoding/json"
	"fmt"
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
		var tx TransactionRequest
		err := ws.lib.DecodeJSONBody(w, req, &tx)
		if !ws.lib.IsHttpOk(err, w) {
			return
		}
		if !tx.Validate() {
			log.Println("ERROR: invalid payload")
			w.WriteHeader(http.StatusInternalServerError)
			_, err = io.WriteString(w, string(ws.lib.JsonStatus("failed: invalid payload")))
			return
		}

		publicKey := ws.lib.PublicKeyFromString(*tx.SenderPublicKey)
		privateKey := ws.lib.PrivateKeyFromString(*tx.SenderPrivateKey, publicKey)
		value, err := strconv.ParseFloat(*tx.SenderSendAmount, 32)
		if err != nil {
			log.Println("ERROR: parse amount failed")
		}
		value32 := float32(value)

		fmt.Println("publicKey:", publicKey)
		fmt.Println("privateKey:", privateKey)
		fmt.Println("value32:", value32)

		transaction := wallet.NewTransaction(
			privateKey,
			publicKey,
			*tx.SenderBlockchainAddress,
			*tx.RecipientBlockchainAddress,
			value32,
		)
		signature := transaction.GenerateSignature()
		signatureStr := signature.String()

		btr := &block.TransactionRequest{
			SenderBlockchainAddress:    tx.SenderBlockchainAddress,
			RecipientBlockchainAddress: tx.RecipientBlockchainAddress,
			SenderPublicKey:            tx.SenderPublicKey,
			Value:                      &value32,
			Signature:                  &signatureStr,
		}

		m, _ := json.Marshal(btr)
		buf := bytes.NewBuffer(m)
		blockchainEndpoint := ws.Gateway() + "/transactions"
		log.Println("INFO: Calling blockchain endpoint:", blockchainEndpoint)
		resp, err := http.Post(blockchainEndpoint, "application/json", buf)

		if err != nil {
			log.Println("ERROR: error calling blockchain", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Add("Content-Type", "application/json")
			io.WriteString(w, string(ws.lib.JsonStatus("failed")))
			return
		}

		if resp.StatusCode == 201 {
			w.Header().Add("Content-Type", "application/json")
			io.WriteString(w, string(ws.lib.JsonStatus("success")))
			log.Println("create transaction success")
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, string(ws.lib.JsonStatus(fmt.Sprintf("failed: %v", resp.StatusCode))))
			log.Println("ERROR: create transaction failed:", resp.StatusCode)
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method")
	}
}

func (ws *WalletServer) WalletAmount(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		blockchainAddress := req.URL.Query().Get("blockchain_address")
		endpoint := fmt.Sprintf("%s/amount", ws.Gateway())
		client := &http.Client{}

		bcsReq, _ := http.NewRequest("GET", endpoint, nil)
		q := bcsReq.URL.Query()
		q.Add("blockchain_address", blockchainAddress)
		bcsReq.URL.RawQuery = q.Encode()
		bcsResp, err := client.Do(bcsReq)
		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(w, string(ws.lib.JsonStatus("fail")))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add(ws.lib.GetApplicationJson())
		if bcsResp.StatusCode == 200 {
			decoder := json.NewDecoder(bcsResp.Body)
			var bar block.AmountResponse
			err := decoder.Decode(&bar)
			if err != nil {
				log.Printf("ERROR: %v", err)
				io.WriteString(w, string(ws.lib.JsonStatus("fail")))
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			m, _ := json.Marshal(struct {
				Message string  `json:"message"`
				Amount  float32 `json:"amount"`
			}{
				Message: "success",
				Amount:  bar.Amount,
			})
			io.WriteString(w, string(m[:]))
		} else {
			log.Printf("ERROR: %v", err)
			io.WriteString(w, string(ws.lib.JsonStatus("fail")))
			w.WriteHeader(http.StatusBadRequest)
		}
	default:
		log.Printf("ERROR: Invalid HTTP method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (ws *WalletServer) Run() {
	http.HandleFunc("/", ws.Index)
	http.HandleFunc("/wallet", ws.Wallet)
	http.HandleFunc("/wallet/amount", ws.WalletAmount)
	http.HandleFunc("/transaction", ws.CreateTransaction)
	log.Printf("Running wallet server on port %v\n", ws.Port())
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(ws.Port())), nil))
}
