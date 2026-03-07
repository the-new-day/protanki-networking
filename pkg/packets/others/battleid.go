package others

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Get the id of the last battle you selected
type LastBattleIdPacket struct {
	packets.BasePacket
}

func NewLastBattleIdPacket() *LastBattleIdPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"battleId",
	}

	var id int32 = packets.LastBattleIdID

	return &LastBattleIdPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.LastBattleIdID, "LastBattleId", func() packets.Packet {
		return NewLastBattleIdPacket()
	})
}
