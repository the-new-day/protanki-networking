package lobby

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Loads detailed battle information
type LoadBattleInfoPacket struct {
	packets.BasePacket
}

func NewLoadBattleInfoPacket() *LoadBattleInfoPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewJsonCodec()),
	}

	attributes := []string{
		"json",
	}

	var id int32 = packets.LoadBattleInfoID

	return &LoadBattleInfoPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.LoadBattleInfoID, "LoadBattleInfo", func() packets.Packet {
		return NewLoadBattleInfoPacket()
	})
}
