package utils

import (
	"fmt"
	"io"
	"slices"
)

type ByteBuffer struct {
	data []byte
}

type NotEnoughBytesToRead struct {
	Need int
	Have int
}

func (err *NotEnoughBytesToRead) Error() string {
	return fmt.Sprintf("need %d bytes, have %d", err.Need, err.Have)
}

func NewByteBuffer() *ByteBuffer {
	return &ByteBuffer{data: make([]byte, 0)}
}

func NewByteBufferFromBytes(bytes []byte) *ByteBuffer {
	bb := NewByteBuffer()
	bb.Write(bytes)
	return bb
}

func (b *ByteBuffer) ReadByte() (byte, error) {
	buf, err := b.ReadBytes(1)
	if err != nil {
		return 0, err
	}
	return buf[0], nil
}

func (b *ByteBuffer) WriteByte(value byte) error {
	b.data = append(b.data, value)
	return nil
}

func (b *ByteBuffer) Write(data []byte) (int, error) {
	b.data = append(b.data, data...)
	return len(data), nil
}

func (b *ByteBuffer) ReadBytes(n int) ([]byte, error) {
	if b.Len() < n {
		return nil, &NotEnoughBytesToRead{n, b.Len()}
	}
	result := b.data[:n]
	b.data = b.data[n:]
	return result, nil
}

func (b *ByteBuffer) Read(buf []byte) (int, error) {
	if len(b.data) == 0 {
		return 0, io.EOF
	}

	n := copy(buf, b.data)
	b.data = b.data[n:]

	if n < len(buf) && len(b.data) == 0 {
		return n, io.EOF
	}
	return n, nil
}

func (b *ByteBuffer) Len() int {
	return len(b.data)
}

func (b *ByteBuffer) Bytes() []byte {
	return slices.Clone(b.data)
}
