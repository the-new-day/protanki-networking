package entry

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Answer the captcha
type AnswerCaptchaPacket struct {
	packets.BasePacket
}

func NewAnswerCaptchaPacket() *AnswerCaptchaPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"type_",
		"imagedata",
	}

	var id int32 = packets.AnswerCaptchaID

	return &AnswerCaptchaPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.AnswerCaptchaID, "AnswerCaptcha", func() packets.Packet {
		return NewAnswerCaptchaPacket()
	})
}
