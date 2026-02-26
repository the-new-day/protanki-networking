package garage

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Load garage
type LoadGaragePacket struct {
	packets.BasePacket
}

func NewLoadGaragePacket() *LoadGaragePacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = -479046431

	return &LoadGaragePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
