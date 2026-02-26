package turrets

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Vulcan starts shooting
type VulcanStartOutPacket struct {
	packets.BasePacket
}

func NewVulcanStartOutPacket() *VulcanStartOutPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"clientTime",
	}

	var id int32 = -520655432

	return &VulcanStartOutPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
