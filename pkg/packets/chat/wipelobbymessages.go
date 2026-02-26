package chat

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Wipes all messages by a user in the lobby
type WipeLobbyMessagesPacket struct {
	packets.BasePacket
}

func NewWipeLobbyMessagesPacket() *WipeLobbyMessagesPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"username",
	}

	var id int32 = 1993050216

	return &WipeLobbyMessagesPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
