package turrets

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/custom"
	"github.com/the-new-day/probogo/pkg/codec/multiple"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Player fires a hammer shot
type HammerShotOutPacket struct {
	packets.BasePacket
}

func NewHammerShotOutPacket() *HammerShotOutPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(complex.NewVector3DCodec()),
		codec.Wrap(multiple.NewVectorCodec(custom.NewTargetPositionCodec(), false)),
	}

	attributes := []string{
		"clientTime",
		"direction",
		"shots",
	}

	var id int32 = -541655881

	return &HammerShotOutPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}


func init() {
    packets.Register(-541655881, "HammerShotOutPacket", func() packets.Packet {
        return NewHammerShotOutPacket()
    })
}
