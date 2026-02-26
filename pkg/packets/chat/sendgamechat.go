package chat

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/internal/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
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

	var id int32 = 945463181

	return &SendGameChatPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
