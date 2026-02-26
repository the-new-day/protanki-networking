package custom

import (
	"github.com/the-new-day/probogo/pkg/codec/complex"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
)

type MissionRewardCodec struct {
	CustomCodec
}

func NewMissionRewardCodec() *MissionRewardCodec {
	customCodec := NewCustomCodec(false)

	AddField(customCodec, "amount", &primitive.IntCodec{})
	AddField(customCodec, "name", complex.NewStringCodec())

	return &MissionRewardCodec{CustomCodec: *customCodec}
}
