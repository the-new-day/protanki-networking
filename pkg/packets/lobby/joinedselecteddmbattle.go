package lobby

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/custom"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Sent when a player joins the currently-selected DM battle.
type JoinedSelectedDmBattlePacket struct {
	packets.BasePacket
}

func NewJoinedSelectedDmBattlePacket() *JoinedSelectedDmBattlePacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(custom.NewBattleInfoUserCodec()),
	}

	attributes := []string{
		"battleID",
		"userInfo",
	}

	var id int32 = -911626491

	return &JoinedSelectedDmBattlePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(-911626491, "JoinedSelectedDmBattlePacket", func() packets.Packet {
		return NewJoinedSelectedDmBattlePacket()
	})
}
