package entry

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Login successful
type LoginSuccessPacket struct {
	packets.BasePacket
}

func NewLoginSuccessPacket() *LoginSuccessPacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = packets.LoginSuccessID

	return &LoginSuccessPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.LoginSuccessID, "LoginSuccess", func() packets.Packet {
		return NewLoginSuccessPacket()
	})
}
