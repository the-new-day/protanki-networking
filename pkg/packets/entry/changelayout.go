package entry

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Changes client layout/mode
type ChangeLayoutPacket struct {
	packets.BasePacket
}

func NewChangeLayoutPacket() *ChangeLayoutPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"layout",
	}

	var id int32 = 1118835050

	return &ChangeLayoutPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}


func init() {
    packets.Register(1118835050, "ChangeLayoutPacket", func() packets.Packet {
        return NewChangeLayoutPacket()
    })
}
