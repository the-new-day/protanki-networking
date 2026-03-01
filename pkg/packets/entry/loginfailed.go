package entry

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Login failed
type LoginFailedPacket struct {
	packets.BasePacket
}

func NewLoginFailedPacket() *LoginFailedPacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = 103812952

	return &LoginFailedPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(103812952, "LoginFailed", func() packets.Packet {
		return NewLoginFailedPacket()
	})
}
