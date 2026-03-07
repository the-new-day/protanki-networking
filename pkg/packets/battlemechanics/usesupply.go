package battlemechanics

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Use a supply in the battle
type UseSupplyPacket struct {
	packets.BasePacket
}

func NewUseSupplyPacket() *UseSupplyPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"supply_id",
	}

	var id int32 = packets.UseSupplyID

	return &UseSupplyPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.UseSupplyID, "UseSupply", func() packets.Packet {
		return NewUseSupplyPacket()
	})
}
