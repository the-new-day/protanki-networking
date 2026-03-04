package network

import (
	"bytes"

	"github.com/the-new-day/probogo/pkg/codec"
	"github.com/the-new-day/probogo/pkg/codec/multiple"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
	"github.com/the-new-day/probogo/pkg/packets"
)

// Contains the keys required to activate packet encryption
type ActivateProtectionPacket struct {
	packets.BasePacket
}

func NewActivateProtectionPacket() *ActivateProtectionPacket {
	codecs := []codec.Codec{
		codec.Wrap(multiple.NewVectorCodec(&primitive.ByteCodec{}, false)),
	}

	attributes := []string{
		"keys",
	}

	var id int32 = packets.ActivateProtectionID

	return &ActivateProtectionPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

// NewActivateProtectionPacketWithKeys creates an instance of ActivateProtectionPacket
// and unwraps it with the provided keys.
func NewActivateProtectionPacketWithKeys(keys []byte) *ActivateProtectionPacket {
	buf := &bytes.Buffer{}

	vectorCodec := multiple.NewVectorCodec(&primitive.ByteCodec{}, false)
	vectorCodec.Encode(keys, buf)

	packet := NewActivateProtectionPacket()
	packet.Unwrap(buf)
	return packet
}

func init() {
	packets.Register(packets.ActivateProtectionID, "ActivateProtection", func() packets.Packet {
		return NewActivateProtectionPacket()
	})
}
