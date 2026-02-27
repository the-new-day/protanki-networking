package entry

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Check if a name is up for registration
type CheckNameAvailabilityPacket struct {
	packets.BasePacket
}

func NewCheckNameAvailabilityPacket() *CheckNameAvailabilityPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"username",
	}

	var id int32 = 1083705823

	return &CheckNameAvailabilityPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}


func init() {
    packets.Register(1083705823, "CheckNameAvailabilityPacket", func() packets.Packet {
        return NewCheckNameAvailabilityPacket()
    })
}
