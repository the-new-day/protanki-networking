package lobby

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
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

	var id int32 = 2092412133

	return &SelectBattlePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(2092412133, "SelectBattlePacket", func() packets.Packet {
		return NewSelectBattlePacket()
	})
}
