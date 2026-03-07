package turrets

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Vulcan stops shooting
type VulcanEndOutPacket struct {
	packets.BasePacket
}

func NewVulcanEndOutPacket() *VulcanEndOutPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"clientTime",
	}

	var id int32 = packets.VulcanEndOutID

	return &VulcanEndOutPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.VulcanEndOutID, "VulcanEndOut", func() packets.Packet {
		return NewVulcanEndOutPacket()
	})
}
