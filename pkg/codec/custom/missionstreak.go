package custom

import (
	"github.com/the-new-day/probogo/pkg/codec/primitive"
)

type MissionStreakCodec struct {
	CustomCodec
}

func NewMissionStreakCodec() *MissionStreakCodec {
	customCodec := NewCustomCodec(false)

	AddField(customCodec, "level", &primitive.IntCodec{})
	AddField(customCodec, "streak", &primitive.IntCodec{})
	AddField(customCodec, "doneToday", &primitive.BoolCodec{})
	AddField(customCodec, "questImgID", &primitive.LongCodec{})
	AddField(customCodec, "rewardImgID", &primitive.LongCodec{})

	return &MissionStreakCodec{CustomCodec: *customCodec}
}
