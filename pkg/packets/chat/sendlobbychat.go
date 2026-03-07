package chat

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Sends a chat message to the lobby
type SendLobbyChatPacket struct {
	packets.BasePacket
}

func NewSendLobbyChatPacket() *SendLobbyChatPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"username",
		"message",
	}

	var id int32 = packets.SendLobbyChatID

	return &SendLobbyChatPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.SendLobbyChatID, "SendLobbyChat", func() packets.Packet {
		return NewSendLobbyChatPacket()
	})
}
