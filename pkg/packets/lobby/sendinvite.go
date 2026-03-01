package lobby

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
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

	var id int32 = -864265623

	return &SendInvitePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(-864265623, "SendInvite", func() packets.Packet {
		return NewSendInvitePacket()
	})
}
