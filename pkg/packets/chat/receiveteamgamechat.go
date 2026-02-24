package chat

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/internal/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// A team message has been received from team battle chat.
type ReceiveTeamGameChatPacket struct {
	packets.BasePacket
}

func NewReceiveTeamGameChatPacket() *ReceiveTeamGameChatPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.IntCodec{}),
	}
	attributes := []string{"username", "message", "team"}

	var id int32 = -449356094

	return &ReceiveTeamGameChatPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
