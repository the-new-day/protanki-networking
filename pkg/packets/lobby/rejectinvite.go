package lobby

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Reject a player
type RejectInvitePacket struct {
	packets.BasePacket
}

func NewRejectInvitePacket() *RejectInvitePacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"username",
	}

	var id int32 = 1152865919

	return &RejectInvitePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
