package battleinfo

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/custom"
	"github.com/the-new-day/protanki-networking/pkg/codec/multiple"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Damage dealt to a tank
type TankDamagePacket struct {
	packets.BasePacket
}

func NewTankDamagePacket() *TankDamagePacket {
	codecs := []codec.Codec{
		codec.Wrap(multiple.NewVectorCodec(custom.NewTankDamageCodec(), false)),
	}

	attributes := []string{
		"damages",
	}

	var id int32 = packets.TankDamageID

	return &TankDamagePacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.TankDamageID, "TankDamage", func() packets.Packet {
		return NewTankDamagePacket()
	})
}
