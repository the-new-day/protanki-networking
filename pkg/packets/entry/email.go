package entry

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
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

	var id int32 = 613462801

	return &EmailPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
