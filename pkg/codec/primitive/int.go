package primitive

import (
	"bytes"
	"encoding/binary"

	"github.com/the-new-day/protanki-networking/pkg/utils"
)

type IntCodec struct{}

func (c *IntCodec) Decode(buf *bytes.Buffer) (int32, error) {
	byteBuf, err := utils.ReadBytes(4, buf)
	if err != nil {
		return 0, err
	}
	return int32(binary.BigEndian.Uint32(byteBuf)), nil
}

func (c *IntCodec) Encode(value int32, buf *bytes.Buffer) (int, error) {
	var byteBuf [4]byte
	binary.BigEndian.PutUint32(byteBuf[:], uint32(value))
	return buf.Write(byteBuf[:])
}
