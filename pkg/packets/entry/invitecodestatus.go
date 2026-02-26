package entry

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Do we need an invite code to access the server?
type InviteCodeStatusPacket struct {
	packets.BasePacket
}

func NewInviteCodeStatusPacket() *InviteCodeStatusPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.BoolCodec{}),
	}

	attributes := []string{
		"inviteEnabled",
	}

	var id int32 = 444933603

	return &InviteCodeStatusPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
