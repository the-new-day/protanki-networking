package battlemechanics

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/packets"
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

	var id int32 = packets.UserSuppliesInfoID

	return &UserSuppliesInfoPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.UserSuppliesInfoID, "UserSuppliesInfo", func() packets.Packet {
		return NewUserSuppliesInfoPacket()
	})
}
