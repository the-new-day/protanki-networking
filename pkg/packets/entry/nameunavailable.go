package entry

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Said name is unavailable for registration with a list of alternative suggested usernames
type NameUnavailablePacket struct {
	packets.BasePacket
}

func NewNameUnavailablePacket() *NameUnavailablePacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewVectorStringCodec()),
	}

	attributes := []string{
		"usernames",
	}

	var id int32 = packets.NameUnavailableID

	return &NameUnavailablePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.NameUnavailableID, "NameUnavailable", func() packets.Packet {
		return NewNameUnavailablePacket()
	})
}
