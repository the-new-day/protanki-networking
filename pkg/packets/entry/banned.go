package entry

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Account banned
type BannedPacket struct {
	packets.BasePacket
}

func NewBannedPacket() *BannedPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"reason",
	}

	var id int32 = packets.BannedID

	return &BannedPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.BannedID, "Banned", func() packets.Packet {
		return NewBannedPacket()
	})
}
