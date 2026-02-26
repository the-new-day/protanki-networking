package lobby

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
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

	var id int32 = -2133657895

	return &JoinedOutsideDmBattlePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
