package lobby

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Loads all current battles
type LoadAllBattlesPacket struct {
	packets.BasePacket
}

func NewLoadAllBattlesPacket() *LoadAllBattlesPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewJsonCodec()),
	}

	attributes := []string{
		"battlesJson",
	}

	var id int32 = 552006706

	return &LoadAllBattlesPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
