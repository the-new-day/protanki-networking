package status

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/custom"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Sets the battle status of the player to the battle ID
type InBattleStatusPacket struct {
	packets.BasePacket
}

func NewInBattleStatusPacket() *InBattleStatusPacket {
	codecs := []codec.Codec{
		codec.Wrap(custom.NewBattleNotifierCodec()),
	}

	attributes := []string{
		"battleNotifier",
	}

	var id int32 = -1895446889

	return &InBattleStatusPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(-1895446889, "InBattleStatus", func() packets.Packet {
		return NewInBattleStatusPacket()
	})
}
