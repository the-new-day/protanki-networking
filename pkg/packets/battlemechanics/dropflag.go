package battlemechanics

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Drop taken flag
type DropFlagPacket struct {
	packets.BasePacket
}

func NewDropFlagPacket() *DropFlagPacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = packets.DropFlagID

	return &DropFlagPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.DropFlagID, "DropFlag", func() packets.Packet {
		return NewDropFlagPacket()
	})
}
