package custom

import (
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/internal/codec/primitive"
)

type BattleInfoCodec struct {
	CustomCodec
}

func NewBattleInfoCodec() *BattleInfoCodec {
	customCodec := NewCustomCodec(false)

	AddField(customCodec, "battleID", &complex.StringCodec{})
	AddField(customCodec, "mapName", &complex.StringCodec{})
	AddField(customCodec, "mode", &primitive.IntCodec{})
	AddField(customCodec, "private", &primitive.BoolCodec{})
	AddField(customCodec, "proBattle", &primitive.BoolCodec{})
	AddField(customCodec, "range", NewRankRangeCodec())
	AddField(customCodec, "serverNumber", &primitive.IntCodec{})

	return &BattleInfoCodec{CustomCodec: *customCodec}
}
