package battlemechanics

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Load Bonus Box Resources
type UserSuppliesInfoPacket struct {
	packets.BasePacket
}

func NewUserSuppliesInfoPacket() *UserSuppliesInfoPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"json",
	}

	var id int32 = -137249251

	return &UserSuppliesInfoPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(-137249251, "UserSuppliesInfo", func() packets.Packet {
		return NewUserSuppliesInfoPacket()
	})
}
