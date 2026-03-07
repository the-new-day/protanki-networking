package primitive

import (
	"bytes"
	"encoding/binary"
	"math"

	"github.com/the-new-day/protanki-networking/pkg/utils"
)

type FloatCodec struct{}

func (c *FloatCodec) Decode(buf *bytes.Buffer) (float32, error) {
	byteBuf, err := utils.ReadBytes(4, buf)
	if err != nil {
		return 0, err
	}
	bits := binary.BigEndian.Uint32(byteBuf)
	return math.Float32frombits(bits), nil
}

func (c *FloatCodec) Encode(value float32, buf *bytes.Buffer) (int, error) {
	var byteBuf [4]byte
	binary.BigEndian.PutUint32(byteBuf[:], math.Float32bits(value))
	return buf.Write(byteBuf[:])
}
