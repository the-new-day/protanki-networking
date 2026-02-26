package custom

import (
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
)

type BattleInfoCodec struct {
	CustomCodec
}

func NewBattleInfoCodec() *BattleInfoCodec {
	customCodec := NewCustomCodec(false)

	AddField(customCodec, "battleID", complex.NewStringCodec())
	AddField(customCodec, "mapName", complex.NewStringCodec())
	AddField(customCodec, "mode", &primitive.IntCodec{})
	AddField(customCodec, "private", &primitive.BoolCodec{})
	AddField(customCodec, "proBattle", &primitive.BoolCodec{})
	AddField(customCodec, "range", NewRankRangeCodec())
	AddField(customCodec, "serverNumber", &primitive.IntCodec{})

	return &BattleInfoCodec{CustomCodec: *customCodec}
}
