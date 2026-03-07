package status

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/packets"
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

	var id int32 = packets.SubscribeStatusID

	return &SubscribeStatusPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.SubscribeStatusID, "SubscribeStatus", func() packets.Packet {
		return NewSubscribeStatusPacket()
	})
}
