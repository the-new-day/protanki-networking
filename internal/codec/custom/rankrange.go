package custom

import (
	"github.com/the-new-day/probogo/internal/codec/complex"
)

type RankRangeCodec struct {
	complex.DoubleIntCodec
}

func NewRankRangeCodec() *RankRangeCodec {
	return &RankRangeCodec{
		DoubleIntCodec: *complex.NewDoubleIntCodec("maxRank", "minRank"),
	}
}
