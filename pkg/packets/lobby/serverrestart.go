package lobby

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
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

	var id int32 = -1712113407

	return &ServerRestartPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
