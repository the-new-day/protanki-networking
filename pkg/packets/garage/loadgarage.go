package garage

import (
	"github.com/the-new-day/probogo/pkg/codec"
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


func init() {
    packets.Register(-479046431, "LoadGaragePacket", func() packets.Packet {
        return NewLoadGaragePacket()
    })
}
