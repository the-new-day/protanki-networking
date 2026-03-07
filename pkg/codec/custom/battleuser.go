package custom

import (
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
)

type BattleUserCodec struct {
	CustomCodec
}

func NewBattleUserCodec() *BattleUserCodec {
	customCodec := NewCustomCodec(false)

	AddField(customCodec, "modLevel", &primitive.IntCodec{})
	AddField(customCodec, "deaths", &primitive.ShortCodec{})
	AddField(customCodec, "kills", &primitive.ShortCodec{})
	AddField(customCodec, "rank", &primitive.ByteCodec{})
	AddField(customCodec, "score", &primitive.IntCodec{})
	AddField(customCodec, "username", complex.NewStringCodec())

	return &BattleUserCodec{CustomCodec: *customCodec}
}
