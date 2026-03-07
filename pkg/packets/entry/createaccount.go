package entry

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
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

	var id int32 = packets.CreateAccountID

	return &CreateAccountPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.CreateAccountID, "CreateAccount", func() packets.Packet {
		return NewCreateAccountPacket()
	})
}
