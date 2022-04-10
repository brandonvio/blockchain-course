package main

import (
	"blockchain/wallet"
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
}

func NewWalletServer(port uint, gateway string) *WalletServer {
	return &WalletServer{
		port:    port,
		gateway: gateway,
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

func (ws *WalletServer) Run() {
	http.HandleFunc("/wallet", ws.Wallet)
	http.HandleFunc("/", ws.Index)
	log.Printf("Running wallet server on port %v\n", ws.Port())
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(ws.Port())), nil))
}
