package block

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Transaction struct {
	sendBlockchainAddress      string
	recipientBlockchainAddress string
	value                      float32
}

func NewTransaction(sender string, recipient string, value float32) *Transaction {
	return &Transaction{
		sendBlockchainAddress:      sender,
		recipientBlockchainAddress: recipient,
		value:                      value,
	}
}

func (t *Transaction) Print() {
	fmt.Printf("%v\n", strings.Repeat("~", 42))
	fmt.Printf("\tsendBlockchainAddress        %s\n", t.sendBlockchainAddress)
	fmt.Printf("\trecipientBlockchainAddress   %s\n", t.recipientBlockchainAddress)
	fmt.Printf("\tvalue                        %.2f\n", t.value)
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"send_blockchain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    t.sendBlockchainAddress,
		Recipient: t.recipientBlockchainAddress,
		Value:     t.value,
	})
}
