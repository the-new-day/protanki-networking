package lobby

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Sent when a player joins a global Team battle, outside from the observer
type JoinedOutsideTeamBattlePacket struct {
	packets.BasePacket
}

func NewJoinedOutsideTeamBattlePacket() *JoinedOutsideTeamBattlePacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"battleID",
		"username",
		"team",
	}

	var id int32 = -169305322

	return &JoinedOutsideTeamBattlePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(-169305322, "JoinedOutsideTeamBattlePacket", func() packets.Packet {
		return NewJoinedOutsideTeamBattlePacket()
	})
}
