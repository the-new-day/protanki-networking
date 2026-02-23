package primitive

import (
	"bytes"
	"encoding/binary"

	"github.com/the-new-day/probogo/internal/utils"
)

type ShortCodec struct{}

func (c *ShortCodec) Decode(buf *bytes.Buffer) (int16, error) {
	byteBuf, err := utils.ReadBytes(2, buf)
	if err != nil {
		return 0, err
	}
	return int16(binary.BigEndian.Uint16(byteBuf)), nil
}

func (c *ShortCodec) Encode(value int16, buf *bytes.Buffer) (int, error) {
	var byteBuf [2]byte
	binary.BigEndian.PutUint16(byteBuf[:], uint16(value))
	return buf.Write(byteBuf[:])
}
