package lobby

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Updates a team
type UpdateTeamBattlePreviewPacket struct {
	packets.BasePacket
}

func NewUpdateTeamBattlePreviewPacket() *UpdateTeamBattlePreviewPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"battleID",
		"team",
		"score",
	}

	var id int32 = 1428217189

	return &UpdateTeamBattlePreviewPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(1428217189, "UpdateTeamBattlePreview", func() packets.Packet {
		return NewUpdateTeamBattlePreviewPacket()
	})
}
