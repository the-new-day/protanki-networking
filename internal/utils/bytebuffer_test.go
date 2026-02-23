package utils

import (
	"io"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewByteBuffer(t *testing.T) {
	buf := NewByteBuffer()
	assert.Equal(t, 0, buf.Len(), "New buffer should be empty")
}

func TestNewByteBufferFromBytes(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03}
	buf := NewByteBufferFromBytes(data)

	assert.Equal(t, 3, buf.Len(), "Buffer length should match input data")

	data[0] = 0xFF
	assert.Equal(t, byte(0x01), buf.data[0], "Buffer should be independent copy")
}

func TestWriteReadInt(t *testing.T) {
	buf := NewByteBuffer()

	values := []int32{0, 42, -42, 2147483647, -2147483648}

	for _, expected := range values {
		buf.WriteInt(expected)
		actual, err := buf.ReadInt()

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
		assert.Equal(t, 0, buf.Len(), "Buffer should be empty after read")
	}
}

func TestWriteReadByte(t *testing.T) {
	buf := NewByteBuffer()

	values := []byte{0, 1, 0xFF, 0x7F, 0x80}

	for _, expected := range values {
		buf.WriteByte(expected)
		actual, err := buf.ReadByte()

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
		assert.Equal(t, 0, buf.Len())
	}
}

func TestWriteReadBool(t *testing.T) {
	buf := NewByteBuffer()

	values := []bool{true, false, true, true, false}

	for _, expected := range values {
		buf.WriteBool(expected)
		actual, err := buf.ReadBool()

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
		assert.Equal(t, 0, buf.Len())
	}
}

func TestWriteReadLong(t *testing.T) {
	buf := NewByteBuffer()

	values := []int64{0, 42, -42, 9223372036854775807, -9223372036854775808}

	for _, expected := range values {
		buf.WriteLong(expected)
		actual, err := buf.ReadLong()

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
		assert.Equal(t, 0, buf.Len())
	}
}

func TestWriteReadShort(t *testing.T) {
	buf := NewByteBuffer()

	values := []int16{0, 42, -42, 32767, -32768}

	for _, expected := range values {
		buf.WriteShort(expected)
		actual, err := buf.ReadShort()

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
		assert.Equal(t, 0, buf.Len())
	}
}

func TestWriteReadFloat(t *testing.T) {
	buf := NewByteBuffer()

	values := []float32{0.0, 3.14, -3.14, 1.5e-10, math.MaxFloat32, math.SmallestNonzeroFloat32}
	epsilon := 0.0001

	for _, expected := range values {
		buf.WriteFloat(expected)
		actual, err := buf.ReadFloat()

		assert.NoError(t, err)
		assert.InDelta(t, expected, actual, epsilon, "Float values should match within delta")
		assert.Equal(t, 0, buf.Len())
	}
}

func TestWriteReadString(t *testing.T) {
	buf := NewByteBuffer()

	values := []string{
		"",
		"Hello",
		"Привет мир!",
		"Special chars: \n\t\r",
		"Very long string " + string(make([]byte, 1000)),
	}

	for _, expected := range values {
		buf.WriteString(expected)
		actual, err := buf.ReadString(len(expected))

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
		assert.Equal(t, 0, buf.Len())
	}
}

func TestMixedReadWrite(t *testing.T) {
	buf := NewByteBuffer()

	// Write different types sequentially
	buf.WriteInt(42)
	buf.WriteBool(true)
	buf.WriteString("Hello")
	buf.WriteFloat(3.14)
	buf.WriteByte(0xFF)

	// Reading in the same order
	i, err := buf.ReadInt()
	assert.NoError(t, err)
	assert.Equal(t, int32(42), i)

	b, err := buf.ReadBool()
	assert.NoError(t, err)
	assert.Equal(t, true, b)

	s, err := buf.ReadString(5)
	assert.NoError(t, err)
	assert.Equal(t, "Hello", s)

	f, err := buf.ReadFloat()
	assert.NoError(t, err)
	assert.InDelta(t, 3.14, f, 0.0001)

	by, err := buf.ReadByte()
	assert.NoError(t, err)
	assert.Equal(t, byte(0xFF), by)

	assert.Equal(t, 0, buf.Len(), "Buffer should be empty after reading all")
}

func TestReadErrors(t *testing.T) {
	buf := NewByteBuffer()

	// Read from empty buffer
	_, err := buf.ReadInt()
	assert.Error(t, err, "ReadInt() on empty buffer should return error")

	// Write one byte, try to read int (needs 4)
	buf.WriteByte(0x01)
	_, err = buf.ReadInt()
	assert.Error(t, err, "ReadInt() with insufficient data should return error")

	// Check error type
	var notEnoughErr *NotEnoughBytesToRead
	assert.ErrorAs(t, err, &notEnoughErr, "Error should be of type NotEnoughBytesToRead")

	// Check error fields
	assert.Equal(t, 4, notEnoughErr.need)
	assert.Equal(t, 1, notEnoughErr.have)
}

func TestWrite(t *testing.T) {
	buf := NewByteBuffer()

	data := []byte{0x01, 0x02, 0x03}
	n, err := buf.Write(data)

	assert.NoError(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, 3, buf.Len())
	assert.Equal(t, data, buf.data)
}

func TestRead_ioReader(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04}
	buf := NewByteBufferFromBytes(data)

	p := make([]byte, 2)

	// First read
	n, err := buf.Read(p)
	assert.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.Equal(t, []byte{0x01, 0x02}, p)

	// Second read
	n, err = buf.Read(p)
	assert.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.Equal(t, []byte{0x03, 0x04}, p)

	// Third read - buffer empty
	n, err = buf.Read(p)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, 0, n)
}

func TestReadBytes(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04}
	buf := NewByteBufferFromBytes(data)

	// Successful read
	readBytes, err := buf.ReadBytes(2)
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x01, 0x02}, readBytes)
	assert.Equal(t, 2, buf.Len())

	// Error - not enough data
	_, err = buf.ReadBytes(3)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "need 3 bytes, have 2")
}

func TestLen(t *testing.T) {
	buf := NewByteBuffer()
	assert.Equal(t, 0, buf.Len())

	buf.WriteString("test")
	assert.Equal(t, 4, buf.Len())

	buf.ReadBytes(2)
	assert.Equal(t, 2, buf.Len())

	buf.ReadBytes(2)
	assert.Equal(t, 0, buf.Len())
}

func TestByteWriterInterface(t *testing.T) {
	buf := NewByteBuffer()

	err := buf.WriteByte(0x42)
	assert.NoError(t, err)
	assert.Equal(t, 1, buf.Len())
	assert.Equal(t, byte(0x42), buf.data[0])
}

func TestWriteBoolEdgeCases(t *testing.T) {
	buf := NewByteBuffer()

	buf.WriteBool(true)
	buf.WriteBool(false)

	assert.Equal(t, 2, buf.Len())

	val1, err := buf.ReadBool()
	assert.NoError(t, err)
	assert.Equal(t, true, val1)

	val2, err := buf.ReadBool()
	assert.NoError(t, err)
	assert.Equal(t, false, val2)

	assert.Equal(t, 0, buf.Len())
}

func TestBigEndian(t *testing.T) {
	buf := NewByteBuffer()

	// 0x12345678 in big-endian should be [0x12, 0x34, 0x56, 0x78]
	buf.WriteInt(0x12345678)

	expected := []byte{0x12, 0x34, 0x56, 0x78}
	assert.Equal(t, expected, buf.data)

	// Read back and verify
	val, err := buf.ReadInt()
	assert.NoError(t, err)
	assert.Equal(t, int32(0x12345678), val)
}

func TestFloatSpecialValues(t *testing.T) {
	buf := NewByteBuffer()

	// Test NaN
	buf.WriteFloat(float32(math.NaN()))
	f, err := buf.ReadFloat()
	assert.NoError(t, err)
	assert.True(t, math.IsNaN(float64(f)))

	// Test +Inf
	buf.WriteFloat(float32(math.Inf(1)))
	f, err = buf.ReadFloat()
	assert.NoError(t, err)
	assert.True(t, math.IsInf(float64(f), 1))

	// Test -Inf
	buf.WriteFloat(float32(math.Inf(-1)))
	f, err = buf.ReadFloat()
	assert.NoError(t, err)
	assert.True(t, math.IsInf(float64(f), -1))
}

func TestMultipleWritesReads(t *testing.T) {
	buf := NewByteBuffer()

	// Write multiple values
	for i := range 10 {
		buf.WriteInt(int32(i))
		buf.WriteBool(i%2 == 0)
	}

	// Read them back
	for i := range 10 {
		val, err := buf.ReadInt()
		assert.NoError(t, err)
		assert.Equal(t, int32(i), val)

		b, err := buf.ReadBool()
		assert.NoError(t, err)
		assert.Equal(t, i%2 == 0, b)
	}

	assert.Equal(t, 0, buf.Len())
}

func TestReadWriteAfterClear(t *testing.T) {
	buf := NewByteBuffer()

	// Write and read some data
	buf.WriteInt(42)
	buf.ReadInt()

	// Buffer should be empty
	assert.Equal(t, 0, buf.Len())

	// Can write again
	buf.WriteInt(100)
	assert.Equal(t, 4, buf.Len())

	val, err := buf.ReadInt()
	assert.NoError(t, err)
	assert.Equal(t, int32(100), val)
}
