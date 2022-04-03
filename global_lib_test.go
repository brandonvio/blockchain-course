package main

import (
	types "blockchain/blockchaintypes"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGlobals_EmptyByte32(t *testing.T) {
	// setup
	globals := NewGlobals()

	Convey("block1 was created as expected", t, func() {
		So(globals.EmptyByte32(), ShouldEqual, types.Byte32{})
		So(globals.NowUnixNano(), ShouldAlmostEqual, time.Now().UnixNano())
	})
}
