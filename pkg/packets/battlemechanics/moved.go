package battlemechanics

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/custom"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Receives movement data of a player from the server.
type MovedPacket struct {
	packets.BasePacket
}

func NewMovedPacket() *MovedPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(custom.NewMoveCodec()),
	}

	attributes := []string{
		"username",
		"movement",
	}

	var id int32 = packets.MovedID

	return &MovedPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.MovedID, "Moved", func() packets.Packet {
		return NewMovedPacket()
	})
}
