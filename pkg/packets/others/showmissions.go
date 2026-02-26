package others

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/custom"
	"github.com/the-new-day/probogo/internal/codec/multiple"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Show the list of missions currently available to the player
type ShowMissionsPacket struct {
	packets.BasePacket
}

func NewShowMissionsPacket() *ShowMissionsPacket {
	codecs := []codec.Codec{
		codec.Wrap(multiple.NewVectorCodec(custom.NewMissionCodec(), false)),
		codec.Wrap(custom.NewMissionStreakCodec()),
	}

	attributes := []string{
		"missions",
		"weeklyStreakInfo",
	}

	var id int32 = 809822533

	return &ShowMissionsPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
