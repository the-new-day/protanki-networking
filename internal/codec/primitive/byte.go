package primitive

import "io"

type ByteCodec struct{}

func (c *ByteCodec) Decode(reader io.Reader) (byte, error) {
	var buf [1]byte
	_, err := reader.Read(buf[:])
	if err != nil {
		return 0, err
	}
	return buf[0], nil
}

func (c *ByteCodec) Encode(value byte, writer io.Writer) (int, error) {
	return writer.Write([]byte{value})
}
