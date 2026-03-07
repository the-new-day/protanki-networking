package lobby

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Removes a battle from the lobby
type RemoveBattlePacket struct {
	packets.BasePacket
}

func NewRemoveBattlePacket() *RemoveBattlePacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"battleID",
	}

	var id int32 = packets.RemoveBattleID

	return &RemoveBattlePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.RemoveBattleID, "RemoveBattle", func() packets.Packet {
		return NewRemoveBattlePacket()
	})
}
