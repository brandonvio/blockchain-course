package main

import "time"

type GlobalLib struct{}

type IGlobalLib interface {
	NowUnixNano() int64
	EmptyByte32() [32]byte
}

func NewGlobals() IGlobalLib {
	return &GlobalLib{}
}

func (g *GlobalLib) NowUnixNano() int64 {
	return time.Now().UnixNano()
}

func (g *GlobalLib) EmptyByte32() [32]byte {
	return [32]byte{}
}
