package chat

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Received a system message in the game chat
type ReceiveGameSystemChatPacket struct {
	packets.BasePacket
}

func NewReceiveGameSystemChatPacket() *ReceiveGameSystemChatPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"message",
	}

	var id int32 = 606668848

	return &ReceiveGameSystemChatPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(606668848, "ReceiveGameSystemChatPacket", func() packets.Packet {
		return NewReceiveGameSystemChatPacket()
	})
}
