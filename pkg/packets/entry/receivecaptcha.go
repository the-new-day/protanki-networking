package entry

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/multiple"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Received a captcha image with its type
type ReceiveCaptchaPacket struct {
	packets.BasePacket
}

func NewReceiveCaptchaPacket() *ReceiveCaptchaPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(multiple.NewVectorCodec(&primitive.ByteCodec{}, false)),
	}

	attributes := []string{
		"type_",
		"imagedata",
	}

	var id int32 = -1670408519

	return &ReceiveCaptchaPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
