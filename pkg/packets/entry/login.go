package entry

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Login information sent by the client
type LoginPacket struct {
	packets.BasePacket
}

func NewLoginPacket() *LoginPacket {
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

	var id int32 = -739684591

	return &LoginPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(-739684591, "LoginPacket", func() packets.Packet {
		return NewLoginPacket()
	})
}
