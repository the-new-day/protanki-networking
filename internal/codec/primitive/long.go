package primitive

import (
	"bytes"
	"encoding/binary"

	"github.com/the-new-day/probogo/internal/utils"
)

type LongCodec struct{}

func (c *LongCodec) Decode(buf *bytes.Buffer) (int64, error) {
	byteBuf, err := utils.ReadBytes(8, buf)
	if err != nil {
		return 0, err
	}
	return int64(binary.BigEndian.Uint64(byteBuf)), nil
}

func (c *LongCodec) Encode(value int64, buf *bytes.Buffer) (int, error) {
	var byteBuf [8]byte
	binary.BigEndian.PutUint64(byteBuf[:], uint64(value))
	return buf.Write(byteBuf[:])
}
