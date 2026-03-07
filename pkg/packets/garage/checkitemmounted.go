package garage

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/packets"
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

	var id int32 = packets.CheckItemMountedID

	return &CheckItemMountedPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.CheckItemMountedID, "CheckItemMounted", func() packets.Packet {
		return NewCheckItemMountedPacket()
	})
}
