package others

import (
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/internal/codec/custom"
	"github.com/the-new-day/probogo/internal/codec/multiple"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Server sends us our referral data
type ReferralDataPacket struct {
	packets.BasePacket
}

func NewReferralDataPacket() *ReferralDataPacket {
	codecs := []codec.Codec{
		codec.Wrap(multiple.NewVectorCodec(custom.NewReferralDataCodec(), false)),
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(complex.NewStringCodec()),
	}

	attributes := []string{
		"referralData",
		"inviteLink",
		"banner",
		"inviteMessage",
	}

	var id int32 = 1587315905

	return &ReferralDataPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}
