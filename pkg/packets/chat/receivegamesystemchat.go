package chat

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// A system message received from the game chat.
type ReceiveGameSystemChatPacket struct {
	packets.BasePacket
}

func NewReceiveGameSystemChatPacket() *ReceiveGameSystemChatPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}
	attributes := []string{"message"}

	var id int32 = 606668848

	return &ReceiveGameSystemChatPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
