package primitive

import (
	"io"

	"github.com/the-new-day/probogo/internal/utils"
)

type BoolCodec struct{}

func (c *BoolCodec) Decode(reader io.Reader) (bool, error) {
	readByte, err := utils.ReadBytes(1, reader)
	if err != nil {
		return false, err
	}
	return readByte[0] != 0, nil
}

func (c *BoolCodec) Encode(value bool, writer io.Writer) (int, error) {
	var v byte
	if value == true {
		v = 1
	} else {
		v = 0
	}
	return writer.Write([]byte{v})
}
