package entry

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Login successful
type LoginSuccessPacket struct {
	packets.BasePacket
}

func NewLoginSuccessPacket() *LoginSuccessPacket {
	codecs := []codec.Codec{}

	attributes := []string{}

	var id int32 = -1923286328

	return &LoginSuccessPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
