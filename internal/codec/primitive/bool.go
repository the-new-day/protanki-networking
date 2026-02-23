package primitive

import (
	"bytes"
)

type BoolCodec struct{}

func (c *BoolCodec) Decode(buf *bytes.Buffer) (bool, error) {
	readByte, err := buf.ReadByte()
	if err != nil {
		return false, err
	}
	return readByte != 0, nil
}

func (c *BoolCodec) Encode(value bool, buf *bytes.Buffer) (int, error) {
	var v byte
	if value == true {
		v = 1
	} else {
		v = 0
	}

	err := buf.WriteByte(v)
	if err != nil {
		return 0, err
	}
	return 1, nil
}
