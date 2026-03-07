package lobby

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
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

	var id int32 = packets.UpdatePlayerTeamBattlePreviewID

	return &UpdatePlayerTeamBattlePreviewPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.UpdatePlayerTeamBattlePreviewID, "UpdatePlayerTeamBattlePreview", func() packets.Packet {
		return NewUpdatePlayerTeamBattlePreviewPacket()
	})
}
