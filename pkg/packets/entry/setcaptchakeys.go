package entry

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/multiple"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Sets captcha hash keys
type SetCaptchaKeysPacket struct {
	packets.BasePacket
}

func NewSetCaptchaKeysPacket() *SetCaptchaKeysPacket {
	codecs := []codec.Codec{
		codec.Wrap(multiple.NewVectorCodec(&primitive.IntCodec{}, false)),
	}

	attributes := []string{
		"keys",
	}

	var id int32 = packets.SetCaptchaKeysID

	return &SetCaptchaKeysPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.SetCaptchaKeysID, "SetCaptchaKeys", func() packets.Packet {
		return NewSetCaptchaKeysPacket()
	})
}
