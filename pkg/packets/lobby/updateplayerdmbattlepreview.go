package lobby

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
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

	var id int32 = packets.UpdatePlayerDmBattlePreviewID

	return &UpdatePlayerDmBattlePreviewPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.UpdatePlayerDmBattlePreviewID, "UpdatePlayerDmBattlePreview", func() packets.Packet {
		return NewUpdatePlayerDmBattlePreviewPacket()
	})
}
