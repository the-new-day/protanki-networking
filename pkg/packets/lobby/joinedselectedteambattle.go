package lobby

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/custom"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Sent when a player joins the currently-selected Team battle.
type JoinedSelectedTeamBattlePacket struct {
	packets.BasePacket
}

func NewJoinedSelectedTeamBattlePacket() *JoinedSelectedTeamBattlePacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(custom.NewBattleInfoUserCodec()),
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"battleID",
		"userInfo",
		"team",
	}

	var id int32 = packets.JoinedSelectedTeamBattleID

	return &JoinedSelectedTeamBattlePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.JoinedSelectedTeamBattleID, "JoinedSelectedTeamBattle", func() packets.Packet {
		return NewJoinedSelectedTeamBattlePacket()
	})
}
