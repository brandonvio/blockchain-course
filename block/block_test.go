package block

import (
	"blockchain/globals"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const BlockTimestamp int64 = 1648402331651366000

func TestBlock_Hash(t *testing.T) {
	// setup
	globals := &globals.GlobalLib{}

	// create blocks
	timestamp, previousHash := BlockTimestamp, globals.EmptyByte32()
	nonce := 7049895176162811509
	transactions := []*Transaction{
		{
			recipientBlockchainAddress: "A",
			senderBlockchainAddress:    "B",
			value:                      100.5,
		},
		{
			recipientBlockchainAddress: "C",
			senderBlockchainAddress:    "D",
			value:                      200.5,
		},
	}

	block := NewBlock(nonce, previousHash, timestamp, transactions)
	block.Print()

	Convey("block1 was created as expected", t, func() {
		So(block.nonce, ShouldEqual, nonce)
		So(block.timestamp, ShouldEqual, timestamp)
		So(block.previousHash, ShouldEqual, previousHash)
		So(fmt.Sprintf("%x", block.Hash()), ShouldEqual, "1af20570d0504ac1322e3cad6f61b3cae60bc31c8131516a00882d50d2a5dfd9")
	})
}
