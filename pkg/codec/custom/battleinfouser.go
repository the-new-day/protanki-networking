package custom

import (
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
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
