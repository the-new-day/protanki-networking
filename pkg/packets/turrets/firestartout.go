package turrets

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Firebird starts shooting
type FireStartOutPacket struct {
	packets.BasePacket
}

func NewFireStartOutPacket() *FireStartOutPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
	}

	attributes := []string{
		"clientTime",
	}

	var id int32 = -1986638927

	return &FireStartOutPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}


func init() {
    packets.Register(-1986638927, "FireStartOutPacket", func() packets.Packet {
        return NewFireStartOutPacket()
    })
}
