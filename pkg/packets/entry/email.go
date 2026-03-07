package entry

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// email
type EmailPacket struct {
	packets.BasePacket
}

func NewEmailPacket() *EmailPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.BoolCodec{}),
	}

	attributes := []string{
		"email",
		"emailConfirmed",
	}

	var id int32 = packets.EmailID

	return &EmailPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.EmailID, "Email", func() packets.Packet {
		return NewEmailPacket()
	})
}
