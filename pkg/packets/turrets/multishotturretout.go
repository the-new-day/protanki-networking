package turrets

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/custom"
	"github.com/the-new-day/protanki-networking/pkg/codec/multiple"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Outgoing shot fired by a multi-shot turret
type MultiShotTurretOutPacket struct {
	packets.BasePacket
}

func NewMultiShotTurretOutPacket() *MultiShotTurretOutPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(complex.NewVector3DCodec()),
		codec.Wrap(multiple.NewVectorCodec(custom.NewTargetPositionCodec(), false)),
	}

	attributes := []string{
		"clientTime",
		"direction",
		"targetHits",
	}

	var id int32 = packets.MultiShotTurretOutID

	return &MultiShotTurretOutPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.MultiShotTurretOutID, "MultiShotTurretOut", func() packets.Packet {
		return NewMultiShotTurretOutPacket()
	})
}
