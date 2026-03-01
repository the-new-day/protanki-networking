package lobby

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Remove this player from the friends count in team status
type LeftOutsideTeamBattlePacket struct {
	packets.BasePacket
}

func NewLeftOutsideTeamBattlePacket() *LeftOutsideTeamBattlePacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"battleID",
		"username",
	}

	var id int32 = 1447204641

	return &LeftOutsideTeamBattlePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(1447204641, "LeftOutsideTeamBattle", func() packets.Packet {
		return NewLeftOutsideTeamBattlePacket()
	})
}
