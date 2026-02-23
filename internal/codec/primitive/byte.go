package primitive

import "bytes"

type ByteCodec struct{}

func (c *ByteCodec) Decode(buf *bytes.Buffer) (byte, error) {
	return buf.ReadByte()
}

func (c *ByteCodec) Encode(value byte, buf *bytes.Buffer) (int, error) {
	err := buf.WriteByte(value)
	if err != nil {
		return 0, err
	}
	return 1, nil
}
