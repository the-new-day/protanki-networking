package lobby

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
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

	var id int32 = 814687528

	return &AcceptInvitePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
