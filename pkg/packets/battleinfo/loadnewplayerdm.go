package battleinfo

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/custom"
	"github.com/the-new-day/probogo/pkg/codec/multiple"
	"github.com/the-new-day/probogo/pkg/packets"
)

// A new player has joined the battle
type LoadNewPlayerDmPacket struct {
	packets.BasePacket
}

func NewLoadNewPlayerDmPacket() *LoadNewPlayerDmPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(multiple.NewVectorCodec(custom.NewBattleUserCodec(), false)),
	}

	attributes := []string{
		"username",
		"userinfos",
	}

	var id int32 = 862913394

	return &LoadNewPlayerDmPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(862913394, "LoadPlayerDm", func() packets.Packet {
		return NewLoadNewPlayerDmPacket()
	})
}
