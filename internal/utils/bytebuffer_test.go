package utils

import (
	"io"
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

func TestWriteReadByte(t *testing.T) {
	buf := NewByteBuffer()

	values := []byte{0, 1, 0xFF, 0x7F, 0x80}

	for _, expected := range values {
		err := buf.WriteByte(expected)
		assert.NoError(t, err)

		actual, err := buf.ReadByte()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	}

	assert.Equal(t, 0, buf.Len(), "Buffer should be empty after all reads")
}

func TestWriteReadBytes(t *testing.T) {
	buf := NewByteBuffer()

	data := []byte{0x01, 0x02, 0x03, 0x04}
	n, err := buf.Write(data)
	assert.NoError(t, err)
	assert.Equal(t, 4, n)
	assert.Equal(t, 4, buf.Len())

	readBytes, err := buf.ReadBytes(2)
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x01, 0x02}, readBytes)
	assert.Equal(t, 2, buf.Len())

	readBytes, err = buf.ReadBytes(2)
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x03, 0x04}, readBytes)
	assert.Equal(t, 0, buf.Len())
}

func TestMixedReadWrite(t *testing.T) {
	buf := NewByteBuffer()

	data1 := []byte{0x01, 0x02}
	data2 := []byte{0x03, 0x04}
	data3 := []byte("Hello")

	buf.Write(data1)
	buf.Write(data2)
	buf.Write(data3)

	read1, err := buf.ReadBytes(2)
	assert.NoError(t, err)
	assert.Equal(t, data1, read1)

	read2, err := buf.ReadBytes(2)
	assert.NoError(t, err)
	assert.Equal(t, data2, read2)

	read3, err := buf.ReadBytes(5)
	assert.NoError(t, err)
	assert.Equal(t, data3, read3)

	assert.Equal(t, 0, buf.Len())
}

func TestReadErrors(t *testing.T) {
	buf := NewByteBuffer()

	_, err := buf.ReadBytes(1)
	assert.Error(t, err)

	var notEnoughErr *NotEnoughBytesToRead
	assert.ErrorAs(t, err, &notEnoughErr)
	assert.Equal(t, 1, notEnoughErr.Need)
	assert.Equal(t, 0, notEnoughErr.Have)

	buf.WriteByte(0x01)
	_, err = buf.ReadBytes(4)
	assert.Error(t, err)
	assert.ErrorAs(t, err, &notEnoughErr)
	assert.Equal(t, 4, notEnoughErr.Need)
	assert.Equal(t, 1, notEnoughErr.Have)
}

func TestWrite(t *testing.T) {
	buf := NewByteBuffer()

	data := []byte{0x01, 0x02, 0x03}
	n, err := buf.Write(data)

	assert.NoError(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, 3, buf.Len())
	assert.Equal(t, data, buf.data)

	assert.Equal(t, len(data), n)
}

func TestRead_ioReader(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04}
	buf := NewByteBufferFromBytes(data)

	p := make([]byte, 2)

	n, err := buf.Read(p)
	assert.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.Equal(t, []byte{0x01, 0x02}, p)
	assert.Equal(t, 2, buf.Len())

	n, err = buf.Read(p)
	assert.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.Equal(t, []byte{0x03, 0x04}, p)
	assert.Equal(t, 0, buf.Len())

	n, err = buf.Read(p)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, 0, n)
}

func TestReadBytes(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04}
	buf := NewByteBufferFromBytes(data)

	readBytes, err := buf.ReadBytes(2)
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x01, 0x02}, readBytes)
	assert.Equal(t, 2, buf.Len())

	_, err = buf.ReadBytes(3)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "need 3 bytes, have 2")
}

func TestLen(t *testing.T) {
	buf := NewByteBuffer()
	assert.Equal(t, 0, buf.Len())

	buf.Write([]byte("test"))
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

func TestMultipleWritesReads(t *testing.T) {
	buf := NewByteBuffer()

	for i := range 10 {
		_, err := buf.Write([]byte{byte(i)})
		assert.NoError(t, err)
	}

	for i := range 10 {
		b, err := buf.ReadBytes(1)
		assert.NoError(t, err)
		assert.Equal(t, byte(i), b[0])
	}

	assert.Equal(t, 0, buf.Len())
}

func TestReadWriteAfterClear(t *testing.T) {
	buf := NewByteBuffer()

	buf.Write([]byte{0x01, 0x02})
	buf.ReadBytes(2)

	assert.Equal(t, 0, buf.Len())

	buf.Write([]byte{0x03, 0x04})
	assert.Equal(t, 2, buf.Len())

	data, err := buf.ReadBytes(2)
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x03, 0x04}, data)
}

func TestBytes(t *testing.T) {
	t.Run("empty buffer", func(t *testing.T) {
		buf := NewByteBuffer()
		bytes := buf.Bytes()

		assert.NotNil(t, bytes)
		assert.Empty(t, bytes)
		assert.Equal(t, 0, len(bytes))
	})

	t.Run("after write", func(t *testing.T) {
		buf := NewByteBuffer()
		buf.Write([]byte{0x01, 0x02, 0x03})

		bytes := buf.Bytes()
		assert.Equal(t, []byte{0x01, 0x02, 0x03}, bytes)
		assert.Equal(t, 3, buf.Len())
	})

	t.Run("returns copy not reference", func(t *testing.T) {
		buf := NewByteBuffer()
		buf.Write([]byte{0x01, 0x02, 0x03})

		bytes := buf.Bytes()
		bytes[0] = 0xFF

		assert.Equal(t, []byte{0x01, 0x02, 0x03}, buf.Bytes())
		assert.Equal(t, byte(0x01), buf.data[0])
	})

	t.Run("multiple calls return different copies", func(t *testing.T) {
		buf := NewByteBuffer()
		buf.Write([]byte{0x01, 0x02, 0x03})

		bytes1 := buf.Bytes()
		bytes2 := buf.Bytes()

		assert.NotSame(t, &bytes1[0], &bytes2[0])
		assert.Equal(t, bytes1, bytes2)
	})

	t.Run("after read", func(t *testing.T) {
		buf := NewByteBuffer()
		buf.Write([]byte{0x01, 0x02, 0x03, 0x04})

		buf.ReadBytes(2)

		bytes := buf.Bytes()
		assert.Equal(t, []byte{0x03, 0x04}, bytes)
	})
}
