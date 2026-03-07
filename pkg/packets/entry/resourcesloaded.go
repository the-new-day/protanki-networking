package entry

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
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

	var id int32 = packets.ResourcesLoadedID

	return &ResourcesLoadedPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.ResourcesLoadedID, "ResourcesLoaded", func() packets.Packet {
		return NewResourcesLoadedPacket()
	})
}
