package custom

import (
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/internal/codec/multiple"
	"github.com/the-new-day/probogo/internal/codec/primitive"
)

type MissionCodec struct {
	CustomCodec
}

func NewMissionCodec() *MissionCodec {
	customCodec := NewCustomCodec(false)

	AddField(customCodec, "freeChange", &primitive.BoolCodec{})
	AddField(customCodec, "description", &complex.StringCodec{})
	AddField(customCodec, "threshold", &primitive.IntCodec{})
	AddField(customCodec, "image", &primitive.LongCodec{})
	AddField(customCodec, "rewards", multiple.NewVectorCodec(&MissionRewardCodec{}, false))
	AddField(customCodec, "progress", &primitive.IntCodec{})
	AddField(customCodec, "missionID", &primitive.IntCodec{})
	AddField(customCodec, "changeCost", &primitive.IntCodec{})

	return &MissionCodec{CustomCodec: *customCodec}
}
