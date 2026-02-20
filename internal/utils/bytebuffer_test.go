package utils

import (
	"bytes"
	"io"
	"math"
	"testing"
)

func TestNewByteBuffer(t *testing.T) {
	buf := NewByteBuffer()
	if buf.Len() != 0 {
		t.Errorf("New buffer should be empty, got length %d", buf.Len())
	}
}

func TestNewByteBufferFromBytes(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03}
	buf := NewByteBufferFromBytes(data)

	if buf.Len() != 3 {
		t.Errorf("Expected length 3, got %d", buf.Len())
	}

	data[0] = 0xFF
	if buf.data[0] != 0x01 {
		t.Errorf("Buffer should be independent copy")
	}
}

func TestWriteReadInt(t *testing.T) {
	buf := NewByteBuffer()

	values := []int32{0, 42, -42, 2147483647, -2147483648}

	for _, expected := range values {
		buf.WriteInt(expected)
		actual, err := buf.ReadInt()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if actual != expected {
			t.Errorf("ReadInt() = %d, want %d", actual, expected)
		}
	}
}

func TestWriteReadByte(t *testing.T) {
	buf := NewByteBuffer()

	values := []byte{0, 1, 0xFF, 0x7F, 0x80}

	for _, expected := range values {
		buf.WriteByte(expected)
		actual, err := buf.ReadByte()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if actual != expected {
			t.Errorf("ReadByte() = %d, want %d", actual, expected)
		}
	}
}

func TestWriteReadBool(t *testing.T) {
	buf := NewByteBuffer()

	values := []bool{true, false, true, true, false}

	for _, expected := range values {
		buf.WriteBool(expected)
		actual, err := buf.ReadBool()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if actual != expected {
			t.Errorf("ReadBool() = %v, want %v", actual, expected)
		}
	}
}

func TestWriteReadLong(t *testing.T) {
	buf := NewByteBuffer()

	values := []int64{0, 42, -42, 9223372036854775807, -9223372036854775808}

	for _, expected := range values {
		buf.WriteLong(expected)
		actual, err := buf.ReadLong()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if actual != expected {
			t.Errorf("ReadLong() = %d, want %d", actual, expected)
		}
	}
}

func TestWriteReadShort(t *testing.T) {
	buf := NewByteBuffer()

	values := []int16{0, 42, -42, 32767, -32768}

	for _, expected := range values {
		buf.WriteShort(expected)
		actual, err := buf.ReadShort()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if actual != expected {
			t.Errorf("ReadShort() = %d, want %d", actual, expected)
		}
	}
}

func TestWriteReadFloat(t *testing.T) {
	buf := NewByteBuffer()

	values := []float32{0.0, 3.14, -3.14, 1.5e-10, math.MaxFloat32, math.SmallestNonzeroFloat32}

	for _, expected := range values {
		buf.WriteFloat(expected)
		actual, err := buf.ReadFloat()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if math.Abs(float64(actual-expected)) > 0.0001 {
			t.Errorf("ReadFloat() = %f, want %f", actual, expected)
		}
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
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if actual != expected {
			t.Errorf("ReadString() = %q, want %q", actual, expected)
		}
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
	i, _ := buf.ReadInt()
	if i != 42 {
		t.Errorf("ReadInt() = %d, want 42", i)
	}

	b, _ := buf.ReadBool()
	if b != true {
		t.Errorf("ReadBool() = %v, want true", b)
	}

	s, _ := buf.ReadString(5)
	if s != "Hello" {
		t.Errorf("ReadString() = %q, want Hello", s)
	}

	f, _ := buf.ReadFloat()
	if math.Abs(float64(f-3.14)) > 0.0001 {
		t.Errorf("ReadFloat() = %f, want 3.14", f)
	}

	by, _ := buf.ReadByte()
	if by != 0xFF {
		t.Errorf("ReadByte() = %d, want 255", by)
	}
}

func TestReadErrors(t *testing.T) {
	buf := NewByteBuffer()

	_, err := buf.ReadInt()
	if err == nil {
		t.Error("ReadInt() on empty buffer should return error")
	}

	buf.WriteByte(0x01)
	_, err = buf.ReadInt()
	if err == nil {
		t.Error("ReadInt() with insufficient data should return error")
	}

	if _, ok := err.(*NotEnoughBytesToRead); !ok {
		t.Errorf("Expected NotEnoughBytesToRead error, got %T", err)
	}
}

func TestWrite(t *testing.T) {
	buf := NewByteBuffer()

	data := []byte{0x01, 0x02, 0x03}
	n, err := buf.Write(data)

	if err != nil {
		t.Errorf("Write() returned error: %v", err)
	}
	if n != 3 {
		t.Errorf("Write() wrote %d bytes, want 3", n)
	}
	if buf.Len() != 3 {
		t.Errorf("Buffer length = %d, want 3", buf.Len())
	}

	if !bytes.Equal(buf.data, data) {
		t.Errorf("Buffer data = %v, want %v", buf.data, data)
	}
}

func TestRead_ioReader(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04}
	buf := NewByteBufferFromBytes(data)

	p := make([]byte, 2)
	n, err := buf.Read(p)

	if err != nil {
		t.Errorf("Read() returned error: %v", err)
	}
	if n != 2 {
		t.Errorf("Read() read %d bytes, want 2", n)
	}
	if !bytes.Equal(p, []byte{0x01, 0x02}) {
		t.Errorf("Read data = %v, want [1 2]", p)
	}

	n, err = buf.Read(p)
	if err != nil {
		t.Errorf("Read() returned error: %v", err)
	}
	if n != 2 {
		t.Errorf("Read() read %d bytes, want 2", n)
	}
	if !bytes.Equal(p, []byte{0x03, 0x04}) {
		t.Errorf("Read data = %v, want [3 4]", p)
	}

	n, err = buf.Read(p)
	if err != io.EOF {
		t.Errorf("Read() on empty buffer should return EOF, got %v", err)
	}
	if n != 0 {
		t.Errorf("Read() on empty buffer should return 0 bytes, got %d", n)
	}
}

func TestReadBytes(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04}
	buf := NewByteBufferFromBytes(data)

	readBytes, err := buf.ReadBytes(2)
	if err != nil {
		t.Errorf("ReadBytes(2) returned error: %v", err)
	}
	if !bytes.Equal(readBytes, []byte{0x01, 0x02}) {
		t.Errorf("ReadBytes(2) = %v, want [1 2]", readBytes)
	}
	if buf.Len() != 2 {
		t.Errorf("Buffer length after read = %d, want 2", buf.Len())
	}

	_, err = buf.ReadBytes(3)
	if err == nil {
		t.Error("ReadBytes(3) should return error when not enough data")
	}
}

func TestLen(t *testing.T) {
	buf := NewByteBuffer()

	if buf.Len() != 0 {
		t.Errorf("Empty buffer length = %d, want 0", buf.Len())
	}

	buf.WriteString("test")
	if buf.Len() != 4 {
		t.Errorf("Buffer length after write = %d, want 4", buf.Len())
	}

	buf.ReadBytes(2)
	if buf.Len() != 2 {
		t.Errorf("Buffer length after read = %d, want 2", buf.Len())
	}
}

func TestByteWriterInterface(t *testing.T) {
	buf := NewByteBuffer()

	err := buf.WriteByte(0x42)
	if err != nil {
		t.Errorf("WriteByte returned error: %v", err)
	}

	if buf.Len() != 1 || buf.data[0] != 0x42 {
		t.Errorf("WriteByte didn't write correctly")
	}
}

func TestWriteBoolEdgeCases(t *testing.T) {
	buf := NewByteBuffer()

	buf.WriteBool(true)
	buf.WriteBool(false)

	if buf.Len() != 2 {
		t.Errorf("Expected 2 bytes, got %d", buf.Len())
	}

	val1, _ := buf.ReadBool()
	val2, _ := buf.ReadBool()

	if val1 != true || val2 != false {
		t.Errorf("WriteBool wrote wrong values: got %v, %v", val1, val2)
	}
}
