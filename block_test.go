package main

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestBlock_Hash(t *testing.T) {
	// setup
	globals := &GlobalLib{}

	// create blocks
	timestamp, previousHash := BlockTimestamp, globals.EmptyByte32()
	nonce := 7049895176162811509
	transactions := []*Transaction{
		{
			recipientBlockchainAddress: "A",
			sendBlockchainAddress:      "B",
			value:                      100.5,
		},
		{
			recipientBlockchainAddress: "C",
			sendBlockchainAddress:      "D",
			value:                      200.5,
		},
	}

	block := NewBlock(nonce, previousHash, timestamp, transactions)
	block.Print()

	Convey("block1 was created as expected", t, func() {
		So(block.nonce, ShouldEqual, nonce)
		So(block.timestamp, ShouldEqual, timestamp)
		So(block.previousHash, ShouldEqual, previousHash)
		So(fmt.Sprintf("%x", block.Hash()), ShouldEqual, "d529bb4e0e2cbb8f66509bd827c2536a19a608911446129860ab38223fa698ed")
	})
}
