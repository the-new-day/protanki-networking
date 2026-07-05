package packets

import (
	"bytes"

	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/modules/protection"
)

// UnknownPacket denotes a package that was not registered
// in the package registry when it was discovered.
type UnknownPacket struct {
	*BasePacket
	data []byte
}

func NewUnknownPacket(id int32, data []byte) *UnknownPacket {
	return &UnknownPacket{
		BasePacket: NewBasePacket(id, []codec.Codec{}, []string{}),
		data:       data,
	}
}

// Unwrap returns the raw packet data in a map with a single key "data".
func (p *UnknownPacket) Unwrap(_ *bytes.Buffer) (map[string]any, error) {
	return map[string]any{"data": p.data}, nil
}

// Wrap encrypts packet data using protection and writes it to a bytes buffer with
// the common packet structure:
//
// [4 bytes length (8 + len(data))] [4 bytes ID] [len(data) bytes data].
func (p *UnknownPacket) Wrap(protection protection.Protection) (*bytes.Buffer, error) {
	data := p.data
	packetLen := HeaderLength + len(data)

	if p.BasePacket.shouldCompress {
		var err error
		data, err = Compress(data)
		if err != nil {
			return nil, err
		}

		packetLen = HeaderLength + len(data)
		packetLen |= 0x40000000 // setting the compression bit
	}

	encrypted := protection.Encrypt(data)

	final := &bytes.Buffer{}
	intCodec := &primitive.IntCodec{}

	intCodec.Encode(int32(packetLen), final)
	intCodec.Encode(p.id, final)

	final.Write(encrypted)
	return final, nil
}
