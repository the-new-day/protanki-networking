package entry

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Information about all maps the client should load
type LoadMapInfoPacket struct {
	packets.BasePacket
}

func NewLoadMapInfoPacket() *LoadMapInfoPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewJsonCodec()),
	}

	attributes := []string{
		"json",
	}

	var id int32 = -838186985

	return &LoadMapInfoPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
