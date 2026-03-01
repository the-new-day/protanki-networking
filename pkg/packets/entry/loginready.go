package entry

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Server sends options for Login
type LoginReadyPacket struct {
	packets.BasePacket
}

func NewLoginReadyPacket() *LoginReadyPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(&primitive.BoolCodec{}),
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"bgResourceID",
		"requireEmail",
		"maxPWLen",
		"minPWLen",
	}

	var id int32 = packets.LoginReadyID

	return &LoginReadyPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.LoginReadyID, "LoginReady", func() packets.Packet {
		return NewLoginReadyPacket()
	})
}
