package others

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Get the id of the last battle you selected
type LastBattleIdPacket struct {
	packets.BasePacket
}

func NewLastBattleIdPacket() *LastBattleIdPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"battleId",
	}

	var id int32 = -602527073

	return &LastBattleIdPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
