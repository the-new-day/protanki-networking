package chat

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Receives a message from the game chat
type ReceiveGameChatPacket struct {
	packets.BasePacket
}

func NewReceiveGameChatPacket() *ReceiveGameChatPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"username",
		"message",
		"team",
	}

	var id int32 = packets.ReceiveGameChatID

	return &ReceiveGameChatPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.ReceiveGameChatID, "ReceiveGameChat", func() packets.Packet {
		return NewReceiveGameChatPacket()
	})
}
