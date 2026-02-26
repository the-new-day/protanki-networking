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

	var id int32 = 810713262

	return &ReceivedInvitePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
