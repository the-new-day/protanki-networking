package status

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Updates a player
type PremiumStatusPacket struct {
	packets.BasePacket
}

func NewPremiumStatusPacket() *PremiumStatusPacket {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"timeLeft",
		"username",
	}

	var id int32 = packets.PremiumStatusID

	return &PremiumStatusPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.PremiumStatusID, "PremiumStatus", func() packets.Packet {
		return NewPremiumStatusPacket()
	})
}
