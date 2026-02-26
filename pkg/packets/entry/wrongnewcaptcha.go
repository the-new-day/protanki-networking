package entry

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/multiple"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// The captcha was incorrect, a new one is sent
type WrongNewCaptchaPacket struct {
	packets.BasePacket
}

func NewWrongNewCaptchaPacket() *WrongNewCaptchaPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(multiple.NewVectorCodec(&primitive.ByteCodec{}, false)),
	}

	attributes := []string{
		"type_",
		"imagedata",
	}

	var id int32 = -373510957

	return &WrongNewCaptchaPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
