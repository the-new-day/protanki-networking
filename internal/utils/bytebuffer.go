package utils

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

type ByteBuffer struct {
	data []byte
}

type NotEnoughBytesToRead struct {
	need int
	have int
}

func (err *NotEnoughBytesToRead) Error() string {
	return fmt.Sprintf("need %d bytes, have %d", err.need, err.have)
}

func NewByteBuffer() *ByteBuffer {
	return &ByteBuffer{data: make([]byte, 0)}
}

func NewByteBufferFromBytes(bytes []byte) *ByteBuffer {
	bb := NewByteBuffer()
	bb.Write(bytes)
	return bb
}

func (b *ByteBuffer) ReadInt() (int32, error) {
	buf, err := b.ReadBytes(4)
	if err != nil {
		return 0, err
	}
	return int32(binary.BigEndian.Uint32(buf)), nil
}

func (b *ByteBuffer) ReadByte() (byte, error) {
	buf, err := b.ReadBytes(1)
	if err != nil {
		return 0, err
	}
	return buf[0], nil
}

func (b *ByteBuffer) ReadBool() (bool, error) {
	readByte, err := b.ReadByte()
	if err != nil {
		return false, err
	}
	return readByte != 0, nil
}

func (b *ByteBuffer) ReadLong() (int64, error) {
	buf, err := b.ReadBytes(8)
	if err != nil {
		return 0, err
	}
	return int64(binary.BigEndian.Uint64(buf)), nil
}

func (b *ByteBuffer) ReadShort() (int16, error) {
	buf, err := b.ReadBytes(2)
	if err != nil {
		return 0, err
	}
	return int16(binary.BigEndian.Uint16(buf)), nil
}

func (b *ByteBuffer) ReadFloat() (float32, error) {
	buf, err := b.ReadBytes(4)
	if err != nil {
		return 0, err
	}
	bits := binary.BigEndian.Uint32(buf)
	return math.Float32frombits(bits), nil
}

func (b *ByteBuffer) ReadString(length int) (string, error) {
	buf, err := b.ReadBytes(length)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func (b *ByteBuffer) WriteInt(value int32) {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, uint32(value))
	b.Write(bytes)
}

func (b *ByteBuffer) WriteByte(value byte) error {
	b.data = append(b.data, value)
	return nil
}

func (b *ByteBuffer) WriteBool(value bool) {
	if value == true {
		b.WriteByte(1)
	} else {
		b.WriteByte(0)
	}
}

func (b *ByteBuffer) WriteLong(value int64) {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(value))
	b.Write(bytes)
}

func (b *ByteBuffer) WriteShort(value int16) {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, uint16(value))
	b.Write(bytes)
}

func (b *ByteBuffer) WriteFloat(value float32) {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, math.Float32bits(value))
	b.Write(bytes)
}

func (b *ByteBuffer) WriteString(value string) {
	b.data = append(b.data, value...)
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
