package primitive

import (
	"encoding/binary"
	"io"

	"github.com/the-new-day/probogo/internal/utils"
)

type LongCodec struct{}

func (c *LongCodec) Decode(reader io.Reader) (int64, error) {
	buf, err := utils.ReadBytes(8, reader)
	if err != nil {
		return 0, err
	}
	return int64(binary.BigEndian.Uint64(buf)), nil
}

func (c *LongCodec) Encode(value int64, writer io.Writer) (int, error) {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(value))
	return writer.Write(bytes)
}
