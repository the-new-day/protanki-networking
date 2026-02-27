package chat

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/custom"
	"github.com/the-new-day/probogo/pkg/codec/multiple"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Receives chat messages from the lobby
type ReceiveLobbyChatPacket struct {
	packets.BasePacket
}

func NewReceiveLobbyChatPacket() *ReceiveLobbyChatPacket {
	codecs := []codec.Codec{
		codec.Wrap(multiple.NewVectorCodec(custom.NewChatMessageCodec(), false)),
	}

	attributes := []string{
		"messages",
	}

	var id int32 = -1263520410

	return &ReceiveLobbyChatPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(-1263520410, "ReceiveLobbyChatPacket", func() packets.Packet {
		return NewReceiveLobbyChatPacket()
	})
}
