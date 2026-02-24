package chat

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// All messages by a user in the lobby has been wiped.
type WipeLobbyChatPacket struct {
	packets.BasePacket
}

func NewWipeLobbyChatPacket() *WipeLobbyChatPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}
	attributes := []string{"username"}

	var id int32 = 1993050216

	return &WipeLobbyChatPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
