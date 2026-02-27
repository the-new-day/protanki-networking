package entry

import (
	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/multiple"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Tells the client which beginner rewards the player has yet to complete
type LoadNewbieRewardsPacket struct {
	packets.BasePacket
}

func NewLoadNewbieRewardsPacket() *LoadNewbieRewardsPacket {
	codecs := []codec.Codec{
		codec.Wrap(multiple.NewVectorCodec(&primitive.IntCodec{}, false)),
	}

	attributes := []string{
		"incompleteRewards",
	}

	var id int32 = 602656160

	return &LoadNewbieRewardsPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(602656160, "LoadbieRewardsPacket", func() packets.Packet {
		return NewLoadNewbieRewardsPacket()
	})
}
