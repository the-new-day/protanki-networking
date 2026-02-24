package battleinfo

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/internal/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Announcement that a gold box will drop soon.
type GoldBoxAnnouncementPacket struct {
	packets.BasePacket
}

func NewGoldBoxAnnouncementPacket() *GoldBoxAnnouncementPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.IntCodec{}),
	}
	attributes := []string{"text", "soundID"}

	var id int32 = -666893269

	return &GoldBoxAnnouncementPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
