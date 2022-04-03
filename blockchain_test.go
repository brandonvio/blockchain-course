package main

import (
	types "blockchain/blockchaintypes"
	"blockchain/mock_main"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

const BlockTimestamp int64 = 1648402331651366000

func TestBlockchain_CreateBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock GlobalLib
	gl := mock_main.NewMockIGlobalLib(ctrl)
	gl.EXPECT().EmptyByte32().AnyTimes().Return(types.Byte32{})
	gl.EXPECT().NowUnixNano().AnyTimes().Return(BlockTimestamp)
	bc := NewBlockchain(gl)

	Convey("blockchain initialized with root block", t, func() {
		rootBlock := bc.chain[0]
		So(rootBlock.nonce, ShouldEqual, 0)
		So(rootBlock.timestamp, ShouldEqual, BlockTimestamp)
		So(len(bc.transactionPool), ShouldEqual, 0)
	})

	Convey("block with transactions is created", t, func() {
		Convey("transactions are added to the blockchain transaction pool", func() {
			bc.AddTransaction("A", "B", 200.2)
			bc.AddTransaction("X", "Y", 395.2)
			So(len(bc.transactionPool), ShouldEqual, 2)
		})
		Convey("copy transaction pool works", func() {
			copiedTransactions := bc.CopyTransactionPool()
			So(copiedTransactions, ShouldNotBeNil)
			So(len(copiedTransactions), ShouldEqual, 2)
		})
		Convey("when block is created, transactions remove from bool added to block", func() {
			b1 := bc.CreateBlock()
			So(len(bc.transactionPool), ShouldEqual, 0) // 0
			So(len(b1.transactions), ShouldEqual, 2)    // 2
			So(b1.nonce, ShouldEqual, 4814)
			So(b1.timestamp, ShouldEqual, BlockTimestamp)
		})
		Convey("create another block", func() {
			bc.AddTransaction("F", "L", 1000.3443)
			bc.AddTransaction("G", "M", 2000.3232)
			b2 := bc.CreateBlock()
			So(b2.nonce, ShouldNotBeNil)
		})
	})
	bc.Print()
}

func TestBlockchain_AddTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	gl := mock_main.NewMockIGlobalLib(ctrl)
	gl.EXPECT().EmptyByte32().AnyTimes().Return(types.Byte32{})
	gl.EXPECT().NowUnixNano().AnyTimes().Return(BlockTimestamp)

	bc := NewBlockchain(gl)
	Convey("transaction was added to blockchain", t, func() {
		bc.AddTransaction("A", "B", 100.5)
		tr := bc.transactionPool[0]
		So(tr.sendBlockchainAddress, ShouldEqual, "A")
		So(tr.recipientBlockchainAddress, ShouldEqual, "B")
		So(tr.value, ShouldEqual, 100.5)
	})
}
