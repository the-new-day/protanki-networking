package battlemechanics

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Drop taken flag
type DropFlagPacket struct {
	packets.BasePacket
}

func NewDropFlagPacket() *DropFlagPacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = -1832611824

	return &DropFlagPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}


func init() {
    packets.Register(-1832611824, "DropFlagPacket", func() packets.Packet {
        return NewDropFlagPacket()
    })
}
