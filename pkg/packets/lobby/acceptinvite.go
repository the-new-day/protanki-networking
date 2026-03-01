package lobby

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Accepts a player
type AcceptInvitePacket struct {
	packets.BasePacket
}

func NewAcceptInvitePacket() *AcceptInvitePacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"username",
	}

	var id int32 = packets.AcceptInviteID

	return &AcceptInvitePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.AcceptInviteID, "AcceptInvite", func() packets.Packet {
		return NewAcceptInvitePacket()
	})
}
