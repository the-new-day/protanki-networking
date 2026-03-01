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

	var id int32 = packets.PingID

	return &PingPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.PingID, "Ping", func() packets.Packet {
		return NewPingPacket()
	})
}
