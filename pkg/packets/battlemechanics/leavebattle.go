package battlemechanics

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Leaves battle to a layout (0 = Lobby, 1 = Garage)
type LeaveBattlePacket struct {
	packets.BasePacket
}

func NewLeaveBattlePacket() *LeaveBattlePacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"layout",
	}

	var id int32 = packets.LeaveBattleID

	return &LeaveBattlePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.LeaveBattleID, "LeaveBattle", func() packets.Packet {
		return NewLeaveBattlePacket()
	})
}
