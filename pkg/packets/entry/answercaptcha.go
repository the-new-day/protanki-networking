package entry

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
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

	var id int32 = 1271163230

	return &AnswerCaptchaPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}


func init() {
    packets.Register(1271163230, "AnswerCaptchaPacket", func() packets.Packet {
        return NewAnswerCaptchaPacket()
    })
}
