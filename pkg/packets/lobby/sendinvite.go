package lobby

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Send a battle invite to a player
type SendInvitePacket struct {
	packets.BasePacket
}

func NewSendInvitePacket() *SendInvitePacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"username",
		"battleID",
	}

	var id int32 = packets.SendInviteID

	return &SendInvitePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.SendInviteID, "SendInvite", func() packets.Packet {
		return NewSendInvitePacket()
	})
}
