package lobby

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Sent when a player joins a global DM battle, outside from the observer
type JoinedOutsideDmBattlePacket struct {
	packets.BasePacket
}

func NewJoinedOutsideDmBattlePacket() *JoinedOutsideDmBattlePacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"battleID",
		"username",
	}

	var id int32 = packets.JoinedOutsideDmBattleID

	return &JoinedOutsideDmBattlePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.JoinedOutsideDmBattleID, "JoinedOutsideDmBattle", func() packets.Packet {
		return NewJoinedOutsideDmBattlePacket()
	})
}
