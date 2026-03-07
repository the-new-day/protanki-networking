package entry

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
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

	var id int32 = packets.InviteCodeStatusID

	return &InviteCodeStatusPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.InviteCodeStatusID, "InviteCodeStatus", func() packets.Packet {
		return NewInviteCodeStatusPacket()
	})
}
