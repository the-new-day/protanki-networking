package garage

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Mount an item in garage
type MountItemPacket struct {
	packets.BasePacket
}

func NewMountItemPacket() *MountItemPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"item_id",
	}

	var id int32 = packets.MountItemID

	return &MountItemPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.MountItemID, "MountItem", func() packets.Packet {
		return NewMountItemPacket()
	})
}
