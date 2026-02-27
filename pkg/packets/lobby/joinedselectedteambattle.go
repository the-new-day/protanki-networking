package lobby

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/custom"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
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

	var id int32 = 118447426

	return &JoinedSelectedTeamBattlePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}


func init() {
    packets.Register(118447426, "JoinedSelectedTeamBattlePacket", func() packets.Packet {
        return NewJoinedSelectedTeamBattlePacket()
    })
}
