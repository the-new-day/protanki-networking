package chat

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// A message has been sent to the lobby chat.
type SendLobbyChatPacket struct {
	packets.BasePacket
}

func NewSendLobbyChatPacket() *SendLobbyChatPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(complex.NewStringCodec()),
	}
	attributes := []string{"username", "message"}

	var id int32 = 705454610

	return &SendLobbyChatPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
