package network

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/multiple"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Contains the keys required to activate packet encryption
type ActivateProtectionPacket struct {
	packets.BasePacket
}

func NewActivateProtectionPacket() *ActivateProtectionPacket {
	codecs := []codec.Codec{
		codec.Wrap(multiple.NewVectorCodec(&primitive.ByteCodec{}, false)),
	}

	attributes := []string{
		"keys",
	}

	var id int32 = 2001736388

	return &ActivateProtectionPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}


func init() {
    packets.Register(2001736388, "ActivateProtectionPacket", func() packets.Packet {
        return NewActivateProtectionPacket()
    })
}
