package entry

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Create new account
type CreateAccountPacket struct {
	packets.BasePacket
}

func NewCreateAccountPacket() *CreateAccountPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.BoolCodec{}),
	}

	attributes := []string{
		"username",
		"password",
		"rememberMe",
	}

	var id int32 = 427083290

	return &CreateAccountPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(427083290, "CreateAccountPacket", func() packets.Packet {
		return NewCreateAccountPacket()
	})
}
