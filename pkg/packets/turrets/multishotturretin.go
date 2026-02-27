package turrets

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/custom"
	"github.com/the-new-day/probogo/pkg/codec/multiple"
	"github.com/the-new-day/probogo/pkg/packets"
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

	var id int32 = -891286317

	return &MultiShotTurretInPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(-891286317, "MultiShotTurretInPacket", func() packets.Packet {
		return NewMultiShotTurretInPacket()
	})
}
