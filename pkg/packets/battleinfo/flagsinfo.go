package battleinfo

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/custom"
	"github.com/the-new-day/probogo/pkg/codec/multiple"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Retrieve information about flagpoles of the current map
type FlagsInfoPacket struct {
	packets.BasePacket
}

func NewFlagsInfoPacket() *FlagsInfoPacket {
	codecs := []codec.Codec{
		codec.Wrap(custom.NewFlagInfoCodec()),
		codec.Wrap(&primitive.LongCodec{}),
		codec.Wrap(&primitive.LongCodec{}),
		codec.Wrap(custom.NewFlagInfoCodec()),
		codec.Wrap(&primitive.LongCodec{}),
		codec.Wrap(&primitive.LongCodec{}),
		codec.Wrap(multiple.NewMultiCodec([]string{"sfx1", "sfx2", "sfx3", "sfx4"}, &primitive.LongCodec{}, false)),
	}

	attributes := []string{
		"blueflag",
		"blueflag_sprite",
		"blueflag_pedestal_model",
		"redflag",
		"redflag_sprite",
		"redflag_pedestal_model",
		"flag_sfx",
	}

	var id int32 = 789790814

	return &FlagsInfoPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
