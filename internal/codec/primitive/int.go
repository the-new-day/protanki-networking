package primitive

import (
	"encoding/binary"
	"io"

	"github.com/the-new-day/probogo/internal/utils"
)

type IntCodec struct{}

func (c *IntCodec) Decode(reader io.Reader) (int32, error) {
	buf, err := utils.ReadBytes(4, reader)
	if err != nil {
		return 0, err
	}
	return int32(binary.BigEndian.Uint32(buf)), nil
}

func (c *IntCodec) Encode(value int32, writer io.Writer) (int, error) {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], uint32(value))
	return writer.Write(b[:])
}
