package lobby

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
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

	var id int32 = packets.UpdateTeamBattlePreviewID

	return &UpdateTeamBattlePreviewPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.UpdateTeamBattlePreviewID, "UpdateTeamBattlePreview", func() packets.Packet {
		return NewUpdateTeamBattlePreviewPacket()
	})
}
