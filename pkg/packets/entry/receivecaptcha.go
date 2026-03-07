package entry

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/multiple"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
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

	var id int32 = packets.ReceiveCaptchaID

	return &ReceiveCaptchaPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.ReceiveCaptchaID, "ReceiveCaptcha", func() packets.Packet {
		return NewReceiveCaptchaPacket()
	})
}
