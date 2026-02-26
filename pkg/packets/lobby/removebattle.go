package lobby

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Removes a battle from the lobby
type RemoveBattlePacket struct {
	packets.BasePacket
}

func NewRemoveBattlePacket() *RemoveBattlePacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"battleID",
	}

	var id int32 = -1848001147

	return &RemoveBattlePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
