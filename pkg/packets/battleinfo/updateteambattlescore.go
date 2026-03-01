package battleinfo

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Update the score of a team within battle
type UpdateTeamBattleScorePacket struct {
	packets.BasePacket
}

func NewUpdateTeamBattleScorePacket() *UpdateTeamBattleScorePacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"team",
		"score",
	}

	var id int32 = packets.UpdateTeamBattleScoreID

	return &UpdateTeamBattleScorePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.UpdateTeamBattleScoreID, "UpdateTeamBattleScore", func() packets.Packet {
		return NewUpdateTeamBattleScorePacket()
	})
}
