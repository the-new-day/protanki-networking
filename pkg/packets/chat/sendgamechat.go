package chat

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Sends a message to the game chat
type SendGameChatPacket struct {
	packets.BasePacket
}

func NewSendGameChatPacket() *SendGameChatPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.BoolCodec{}),
	}

	attributes := []string{
		"message",
		"teamOnly",
	}

	var id int32 = packets.SendGameChatID

	return &SendGameChatPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.SendGameChatID, "SendGameChat", func() packets.Packet {
		return NewSendGameChatPacket()
	})
}
