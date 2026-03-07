package lobby

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
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

	var id int32 = packets.JoinedOutsideTeamBattleID

	return &JoinedOutsideTeamBattlePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.JoinedOutsideTeamBattleID, "JoinedOutsideTeamBattle", func() packets.Packet {
		return NewJoinedOutsideTeamBattlePacket()
	})
}
