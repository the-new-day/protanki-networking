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

	var id int32 = 1484572481

	return &PongPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
