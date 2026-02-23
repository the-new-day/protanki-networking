package primitive

import (
	"encoding/binary"
	"io"
	"math"

	"github.com/the-new-day/probogo/internal/utils"
)

type FloatCodec struct{}

func (c *FloatCodec) Decode(reader io.Reader) (float32, error) {
	buf, err := utils.ReadBytes(4, reader)
	if err != nil {
		return 0, err
	}
	bits := binary.BigEndian.Uint32(buf)
	return math.Float32frombits(bits), nil
}

func (c *FloatCodec) Encode(value float32, writer io.Writer) (int, error) {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, math.Float32bits(value))
	return writer.Write(bytes)
}
