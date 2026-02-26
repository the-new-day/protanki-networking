package lobby

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/internal/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Updates a player
type UpdatePlayerTeamBattlePreviewPacket struct {
	packets.BasePacket
}

func NewUpdatePlayerTeamBattlePreviewPacket() *UpdatePlayerTeamBattlePreviewPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"battleID",
		"username",
		"score",
	}

	var id int32 = -375282889

	return &UpdatePlayerTeamBattlePreviewPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
