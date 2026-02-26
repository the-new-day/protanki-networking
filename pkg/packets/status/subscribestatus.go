package status

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Subscribe to status updates of a player.
type SubscribeStatusPacket struct {
	packets.BasePacket
}

func NewSubscribeStatusPacket() *SubscribeStatusPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"username",
	}

	var id int32 = 1774907609

	return &SubscribeStatusPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
