package turrets

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
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

	var id int32 = 1794372798

	return &VulcanEndOutPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
