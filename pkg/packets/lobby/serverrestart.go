package lobby

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Indicates time (in seconds) until server restart
type ServerRestartPacket struct {
	packets.BasePacket
}

func NewServerRestartPacket() *ServerRestartPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"time",
	}

	var id int32 = packets.ServerRestartID

	return &ServerRestartPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.ServerRestartID, "ServerRestart", func() packets.Packet {
		return NewServerRestartPacket()
	})
}
