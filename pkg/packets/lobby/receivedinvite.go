package lobby

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Accepts a player
type ReceivedInvitePacket struct {
	packets.BasePacket
}

func NewReceivedInvitePacket() *ReceivedInvitePacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"username",
	}

	var id int32 = packets.ReceivedInviteID

	return &ReceivedInvitePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.ReceivedInviteID, "ReceivedInvite", func() packets.Packet {
		return NewReceivedInvitePacket()
	})
}
