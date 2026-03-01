package entry

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Server tells us what resources to load
type LoadResourcesPacket struct {
	packets.BasePacket
}

func NewLoadResourcesPacket() *LoadResourcesPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewJsonCodec()),
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"json",
		"callbackID",
	}

	var id int32 = packets.LoadResourcesID

	return &LoadResourcesPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.LoadResourcesID, "LoadResources", func() packets.Packet {
		return NewLoadResourcesPacket()
	})
}
