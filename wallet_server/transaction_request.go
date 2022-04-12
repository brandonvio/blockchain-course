package main

type TransactionRequest struct {
	SenderPrivateKey           *string `json:"sender_private_key"`
	SenderBlockchainAddress    *string `json:"sender_blockchain_address"`
	RecipientBlockchainAddress *string `json:"sender_recipient_blockchain_address"`
	SenderPublicKey            *string `json:"sender_public_key"`
	SenderSendAmount           *string `json:"sender_send_amount"`
}

func (tr *TransactionRequest) Validate() bool {
	if tr.SenderPrivateKey == nil ||
		tr.SenderBlockchainAddress == nil ||
		tr.RecipientBlockchainAddress == nil ||
		tr.SenderPublicKey == nil ||
		tr.SenderSendAmount == nil ||
		*tr.SenderPrivateKey == "" ||
		*tr.SenderBlockchainAddress == "" ||
		*tr.RecipientBlockchainAddress == "" ||
		*tr.SenderPublicKey == "" ||
		*tr.SenderSendAmount == "" {
		return false
	}
	return true
}
