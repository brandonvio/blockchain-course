package globals

import (
	types "blockchain/blockchaintypes"
	"time"
)

type GlobalLib struct{}

type IGlobalLib interface {
	NowUnixNano() int64
	EmptyByte32() types.Byte32
}

func NewGlobals() IGlobalLib {
	return &GlobalLib{}
}

func (g *GlobalLib) NowUnixNano() int64 {
	return time.Now().UnixNano()
}

func (g *GlobalLib) EmptyByte32() types.Byte32 {
	return types.Byte32{}
}
