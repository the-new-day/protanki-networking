package garage

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Load garage
type LoadGaragePacket struct {
	packets.BasePacket
}

func NewLoadGaragePacket() *LoadGaragePacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = packets.LoadGarageID

	return &LoadGaragePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.LoadGarageID, "LoadGarage", func() packets.Packet {
		return NewLoadGaragePacket()
	})
}
