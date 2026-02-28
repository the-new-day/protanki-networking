package packets

import (
	"bytes"
	"encoding/binary"

	"github.com/the-new-day/probogo/pkg/modules/protection"
)

// UnknownPacket denotes a package that was not registered
// in the package registry when it was discovered.
type UnknownPacket struct {
	id   int32
	data []byte
}

func NewUnknownPacket(id int32, data []byte) *UnknownPacket {
	return &UnknownPacket{id, data}
}

func (p *UnknownPacket) ID() int32 {
	return p.id
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
	payload := &bytes.Buffer{}
	payload.Write(p.data)
	encrypted := protection.Encrypt(payload.Bytes())

	final := &bytes.Buffer{}

	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, uint32(8+len(encrypted)))
	final.Write(lenBuf)

	idBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(idBuf, uint32(p.id))
	final.Write(idBuf)

	final.Write(encrypted)
	return final, nil
}
