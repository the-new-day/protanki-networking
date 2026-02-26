package network

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Ping Packet from server
type PingPacket struct {
	packets.BasePacket
}

func NewPingPacket() *PingPacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = -555602629

	return &PingPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
