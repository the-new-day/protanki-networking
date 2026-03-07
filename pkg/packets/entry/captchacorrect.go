package entry

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Captcha is correct
type CaptchaCorrectPacket struct {
	packets.BasePacket
}

func NewCaptchaCorrectPacket() *CaptchaCorrectPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"type_",
	}

	var id int32 = packets.CaptchaCorrectID

	return &CaptchaCorrectPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.CaptchaCorrectID, "CaptchaCorrect", func() packets.Packet {
		return NewCaptchaCorrectPacket()
	})
}
