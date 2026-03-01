package lobby

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Updates a player
type UpdatePlayerDmBattlePreviewPacket struct {
	packets.BasePacket
}

func NewUpdatePlayerDmBattlePreviewPacket() *UpdatePlayerDmBattlePreviewPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"battleID",
		"username",
		"kills",
	}

	var id int32 = -1263036614

	return &UpdatePlayerDmBattlePreviewPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(-1263036614, "UpdatePlayerDmBattlePreview", func() packets.Packet {
		return NewUpdatePlayerDmBattlePreviewPacket()
	})
}
