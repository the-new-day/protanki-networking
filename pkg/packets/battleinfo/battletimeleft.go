package battleinfo

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Battle time left in seconds.
type BattleTimeLeftPacket struct {
	packets.BasePacket
}

func NewBattleTimeLeftPacket() *BattleTimeLeftPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
	}
	attributes := []string{"timeLimitInSec"}

	var id int32 = 732434644

	return &BattleTimeLeftPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
