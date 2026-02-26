package garage

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
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

	var id int32 = -1505530736

	return &MountItemPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
