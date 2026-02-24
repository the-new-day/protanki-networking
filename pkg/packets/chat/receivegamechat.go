package chat

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/internal/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// A message received from the game chat.
type ReceiveGameChatPacket struct {
	packets.BasePacket
}

func NewReceiveGameChatPacket() *ReceiveGameChatPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.IntCodec{}),
	}
	attributes := []string{"username", "message", "team"}

	var id int32 = 1259981343

	return &ReceiveGameChatPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
