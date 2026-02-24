package battleinfo

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/internal/codec/custom"
	"github.com/the-new-day/probogo/internal/codec/multiple"
	"github.com/the-new-day/probogo/internal/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// A new player has joined the team battle, the observer being within the battle.
type LoadNewPlayerTeamBattlePacket struct {
	packets.BasePacket
}

func NewLoadNewPlayerTeamBattlePacket() *LoadNewPlayerTeamBattlePacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(multiple.NewVectorCodec(custom.NewBattleUserCodec(), false)),
		codec.Wrap(&primitive.IntCodec{}),
	}
	attributes := []string{
		"username",
		"userinfos",
		"team",
	}

	var id int32 = 2040021062

	return &LoadNewPlayerTeamBattlePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
