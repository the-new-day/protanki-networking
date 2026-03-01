package entry

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Client callsback after finishes resource loading (NOT IMPLEMENTED)
type ResourcesLoadedPacket struct {
	packets.BasePacket
}

func NewResourcesLoadedPacket() *ResourcesLoadedPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"callbackId",
	}

	var id int32 = -82304134

	return &ResourcesLoadedPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(-82304134, "ResourcesLoaded", func() packets.Packet {
		return NewResourcesLoadedPacket()
	})
}
