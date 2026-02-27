package entry

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/packets"
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

	var id int32 = -600078553

	return &BannedPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(-600078553, "BannedPacket", func() packets.Packet {
		return NewBannedPacket()
	})
}
