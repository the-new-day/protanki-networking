package entry

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Login failed
type LoginFailedPacket struct {
	packets.BasePacket
}

func NewLoginFailedPacket() *LoginFailedPacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = packets.LoginFailedID

	return &LoginFailedPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.LoginFailedID, "LoginFailed", func() packets.Packet {
		return NewLoginFailedPacket()
	})
}
