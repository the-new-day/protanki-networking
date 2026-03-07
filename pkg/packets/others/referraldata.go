package others

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/codec/custom"
	"github.com/the-new-day/protanki-networking/pkg/codec/multiple"
	"github.com/the-new-day/protanki-networking/pkg/packets"
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

	var id int32 = packets.ReferralDataID

	return &ReferralDataPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.ReferralDataID, "ReferralData", func() packets.Packet {
		return NewReferralDataPacket()
	})
}
