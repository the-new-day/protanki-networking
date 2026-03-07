package turrets

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/custom"
	"github.com/the-new-day/protanki-networking/pkg/codec/multiple"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Incoming shot fired by a multi-shot turret
type MultiShotTurretInPacket struct {
	packets.BasePacket
}

func NewMultiShotTurretInPacket() *MultiShotTurretInPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(complex.NewVector3DCodec()),
		codec.Wrap(multiple.NewVectorCodec(custom.NewTargetHitCodec(), false)),
	}

	attributes := []string{
		"shooter",
		"direction",
		"targetHits",
	}

	var id int32 = packets.MultiShotTurretInID

	return &MultiShotTurretInPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.MultiShotTurretInID, "MultiShotTurretIn", func() packets.Packet {
		return NewMultiShotTurretInPacket()
	})
}
