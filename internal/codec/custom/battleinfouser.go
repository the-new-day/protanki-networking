package custom

import (
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/internal/codec/primitive"
)

type BattleInfoUserCodec struct {
	CustomCodec
}

func NewBattleInfoUserCodec() *BattleInfoUserCodec {
	customCodec := NewCustomCodec(false)

	AddField(customCodec, "kills", &primitive.IntCodec{})
	AddField(customCodec, "score", &primitive.IntCodec{})
	AddField(customCodec, "suspicious", &primitive.BoolCodec{})
	AddField(customCodec, "user", complex.NewStringCodec())

	return &BattleInfoUserCodec{CustomCodec: *customCodec}
}
