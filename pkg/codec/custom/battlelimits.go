package custom

import (
	"github.com/the-new-day/probogo/pkg/codec/complex"
)

type BattleLimitsCodec struct {
	complex.DoubleIntCodec
}

func NewBattleLimitsCodec() *BattleLimitsCodec {
	return &BattleLimitsCodec{
		DoubleIntCodec: *complex.NewDoubleIntCodec("scoreLimit", "timeLimit"),
	}
}
