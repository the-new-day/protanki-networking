package lobby

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Client selects a battle / Server confirms selection of battle
type SelectBattlePacket struct {
	packets.BasePacket
}

func NewSelectBattlePacket() *SelectBattlePacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"battleID",
	}

	var id int32 = packets.SelectBattleID

	return &SelectBattlePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.SelectBattleID, "SelectBattle", func() packets.Packet {
		return NewSelectBattlePacket()
	})
}
