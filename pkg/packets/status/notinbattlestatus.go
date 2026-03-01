package status

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Sets the player
type NotInBattleStatusPacket struct {
	packets.BasePacket
}

func NewNotInBattleStatusPacket() *NotInBattleStatusPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"username",
	}

	var id int32 = packets.NotInBattleStatusID

	return &NotInBattleStatusPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.NotInBattleStatusID, "NotInBattleStatus", func() packets.Packet {
		return NewNotInBattleStatusPacket()
	})
}
