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

	var id int32 = packets.CheckNameAvailabilityID

	return &CheckNameAvailabilityPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.CheckNameAvailabilityID, "CheckNameAvailability", func() packets.Packet {
		return NewCheckNameAvailabilityPacket()
	})
}
