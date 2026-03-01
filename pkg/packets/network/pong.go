package network

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Pong Packet from client
type PongPacket struct {
	packets.BasePacket
}

func NewPongPacket() *PongPacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = packets.PongID

	return &PongPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.PongID, "Pong", func() packets.Packet {
		return NewPongPacket()
	})
}
