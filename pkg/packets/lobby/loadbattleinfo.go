package lobby

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Loads detailed battle information
type LoadBattleInfoPacket struct {
	packets.BasePacket
}

func NewLoadBattleInfoPacket() *LoadBattleInfoPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewJsonCodec()),
	}

	attributes := []string{
		"json",
	}

	var id int32 = 546722394

	return &LoadBattleInfoPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
