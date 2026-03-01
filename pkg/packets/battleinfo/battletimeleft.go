package battleinfo

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Battle time left in seconds
type BattleTimeLeftPacket struct {
	packets.BasePacket
}

func NewBattleTimeLeftPacket() *BattleTimeLeftPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"timeLimitInSec",
	}

	var id int32 = packets.BattleTimeLeftID

	return &BattleTimeLeftPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.BattleTimeLeftID, "BattleTimeLeft", func() packets.Packet {
		return NewBattleTimeLeftPacket()
	})
}
