package lobby

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Load battle info
type LoadExternalEntrancePacket struct {
	packets.BasePacket
}

func NewLoadExternalEntrancePacket() *LoadExternalEntrancePacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"authUrl",
		"SnId",
	}

	var id int32 = packets.LoadExternalEntranceID

	return &LoadExternalEntrancePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.LoadExternalEntranceID, "LoadExternalEntrance", func() packets.Packet {
		return NewLoadExternalEntrancePacket()
	})
}
