package turrets

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Sends server details about a released railgun shot
type RailgunShotOutPacket struct {
	packets.BasePacket
}

func NewRailgunShotOutPacket() *RailgunShotOutPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(complex.NewVector3DCodec()),
		codec.Wrap(complex.NewVectorStringCodec()),
		codec.Wrap(complex.NewVectorVector3DCodec()),
		codec.Wrap(complex.NewVectorShortCodec()),
		codec.Wrap(complex.NewVectorVector3DCodec()),
		codec.Wrap(complex.NewVectorVector3DCodec()),
	}

	attributes := []string{
		"clientTime",
		"staticHitPoint",
		"targets",
		"targetHitPoints",
		"incarnationIDs",
		"targetBodyPositions",
		"globalHitPoints",
	}

	var id int32 = packets.RailgunShotOutID

	return &RailgunShotOutPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.RailgunShotOutID, "RailgunShotOut", func() packets.Packet {
		return NewRailgunShotOutPacket()
	})
}
