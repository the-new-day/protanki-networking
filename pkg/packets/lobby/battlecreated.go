package lobby

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Loads limited info about a newly created battle
type BattleCreatedPacket struct {
	packets.BasePacket
}

func NewBattleCreatedPacket() *BattleCreatedPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewJsonCodec()),
	}

	attributes := []string{
		"json",
	}

	var id int32 = packets.BattleCreatedID

	return &BattleCreatedPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.BattleCreatedID, "BattleCreated", func() packets.Packet {
		return NewBattleCreatedPacket()
	})
}
