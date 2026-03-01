package lobby

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Censors invalid battle names
type CheckBattleNamePacket struct {
	packets.BasePacket
}

func NewCheckBattleNamePacket() *CheckBattleNamePacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"battleName",
	}

	var id int32 = packets.CheckBattleNameID

	return &CheckBattleNamePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.CheckBattleNameID, "CheckBattleName", func() packets.Packet {
		return NewCheckBattleNamePacket()
	})
}
