package primitive

import (
	"encoding/binary"
	"io"

	"github.com/the-new-day/probogo/internal/utils"
)

type ShortCodec struct{}

func (c *ShortCodec) Decode(reader io.Reader) (int16, error) {
	buf, err := utils.ReadBytes(2, reader)
	if err != nil {
		return 0, err
	}
	return int16(binary.BigEndian.Uint16(buf)), nil
}

func (c *ShortCodec) Encode(value int16, writer io.Writer) (int, error) {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, uint16(value))
	return writer.Write(bytes)
}
