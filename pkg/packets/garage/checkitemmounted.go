package garage

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/internal/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// If the mount was successful or not
type CheckItemMountedPacket struct {
	packets.BasePacket
}

func NewCheckItemMountedPacket() *CheckItemMountedPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.BoolCodec{}),
	}

	attributes := []string{
		"item_id",
		"mounted",
	}

	var id int32 = 2062201643

	return &CheckItemMountedPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
