package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestGlobals_EmptyByte32(t *testing.T) {
	// setup
	globals := NewGlobals()

	Convey("block1 was created as expected", t, func() {
		So(globals.EmptyByte32(), ShouldEqual, [32]byte{})
		So(globals.NowUnixNano(), ShouldAlmostEqual, time.Now().UnixNano())
	})
}
