package status

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Updates Player
type OnlineStatusPacket struct {
	packets.BasePacket
}

func NewOnlineStatusPacket() *OnlineStatusPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.BoolCodec{}),
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"online",
		"serverID",
		"username",
	}

	var id int32 = packets.OnlineStatusID

	return &OnlineStatusPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.OnlineStatusID, "OnlineStatus", func() packets.Packet {
		return NewOnlineStatusPacket()
	})
}
