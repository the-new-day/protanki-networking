package multiple

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/the-new-day/probogo/pkg/codec/primitive"
)

// mockTypedCodec implements codec.TypedCodec for testing
type mockTypedCodec[T any] struct {
	decodeVal    T
	decodeErr    error
	encodeErr    error
	encodeCalled bool
	decodeCalled bool
	lastEncoded  T
	consumeBytes int // how much bytes to "read" in Decode()
}

func (m *mockTypedCodec[T]) Decode(buf *bytes.Buffer) (T, error) {
	m.decodeCalled = true
	if m.decodeErr != nil {
		var zero T
		return zero, m.decodeErr
	}

	if buf.Len() < m.consumeBytes {
		var zero T
		return zero, fmt.Errorf("mock: insufficient data")
	}
	buf.Next(m.consumeBytes)
	return m.decodeVal, nil
}

func (m *mockTypedCodec[T]) Encode(value T, buf *bytes.Buffer) (int, error) {
	m.encodeCalled = true
	m.lastEncoded = value
	if m.encodeErr != nil {
		return 0, m.encodeErr
	}
	return buf.Write([]byte{0xAA}) // Write a marker byte
}

// --- VectorCodec Tests ---

func TestVectorCodec_WithoutBoolshortern(t *testing.T) {
	mockCodec := &mockTypedCodec[int32]{decodeVal: 42}
	codec := NewVectorCodec(mockCodec, false)
	buf := &bytes.Buffer{}

	values := []int32{1, 2, 3, 4, 5}
	n, err := codec.Encode(values, buf)
	assert.NoError(t, err)
	assert.Equal(t, 4+5*1, n)

	result, err := codec.Decode(buf)
	assert.NoError(t, err)
	assert.Len(t, result, 5)

	for _, v := range result {
		assert.Equal(t, int32(42), v)
	}
}

func TestVectorCodec_WithBoolshortern_NotEmpty(t *testing.T) {
	mockCodec := &mockTypedCodec[int32]{decodeVal: 42}
	codec := NewVectorCodec(mockCodec, true)
	buf := &bytes.Buffer{}

	values := []int32{1, 2, 3}
	n, err := codec.Encode(values, buf)
	assert.NoError(t, err)
	assert.Equal(t, 1+4+3*1, n) // empty flag(1) + length(4) + 3 elements

	// First byte should be 0 (not empty)
	assert.Equal(t, byte(0x00), buf.Bytes()[0])

	result, err := codec.Decode(buf)
	assert.NoError(t, err)
	assert.Len(t, result, 3)
}

func TestVectorCodec_WithBoolshortern_Empty(t *testing.T) {
	mockCodec := &mockTypedCodec[int32]{}
	codec := NewVectorCodec(mockCodec, true)
	buf := &bytes.Buffer{}

	values := []int32{}
	n, err := codec.Encode(values, buf)
	assert.NoError(t, err)
	assert.Equal(t, 1, n) // only empty flag

	// First byte should be 1 (empty)
	assert.Equal(t, byte(0x01), buf.Bytes()[0])

	result, err := codec.Decode(buf)
	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, 0, buf.Len())
}

func TestVectorCodec_ErrorHandling(t *testing.T) {
	t.Run("encode with failing codec", func(t *testing.T) {
		mockCodec := &mockTypedCodec[int32]{encodeErr: assert.AnError}
		codec := NewVectorCodec(mockCodec, false)
		buf := &bytes.Buffer{}

		_, err := codec.Encode([]int32{1, 2}, buf)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to encode element")
	})

	t.Run("decode with insufficient data", func(t *testing.T) {
		mockCodec := &mockTypedCodec[int32]{consumeBytes: 1}
		codec := NewVectorCodec(mockCodec, false)
		buf := &bytes.Buffer{}

		// Write only length, not enough data for elements
		intCodec := &primitive.IntCodec{}
		intCodec.Encode(int32(5), buf)

		_, err := codec.Decode(buf)
		assert.Error(t, err)
	})

	t.Run("decode with failing element codec", func(t *testing.T) {
		mockCodec := &mockTypedCodec[int32]{decodeErr: assert.AnError}
		codec := NewVectorCodec(mockCodec, false)
		buf := &bytes.Buffer{}

		intCodec := &primitive.IntCodec{}
		intCodec.Encode(int32(1), buf)
		buf.WriteByte(0xAA) // element data

		_, err := codec.Decode(buf)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decode element")
	})
}

func TestVectorCodec_WithBool(t *testing.T) {
	boolCodec := &primitive.BoolCodec{}
	codec := NewVectorCodec(boolCodec, false)
	buf := &bytes.Buffer{}

	values := []bool{true, false, true}
	n, err := codec.Encode(values, buf)
	assert.NoError(t, err)

	assert.Equal(t, 4+3, n, "should be length(4) + 3 bools")

	// Проверяем содержимое
	result, err := codec.Decode(buf)
	assert.NoError(t, err)
	assert.Equal(t, values, result)
}

func TestVectorCodec_WithBoolAndBoolshortern(t *testing.T) {
	boolCodec := &primitive.BoolCodec{}
	codec := NewVectorCodec(boolCodec, true)
	buf := &bytes.Buffer{}

	values := []bool{true, false, true}
	n, err := codec.Encode(values, buf)
	assert.NoError(t, err)

	assert.Equal(t, 1+4+3, n, "should be flag(1) + length(4) + 3 bools")
}

// --- MultiCodec Tests ---

func TestMultiCodec_WithoutBoolshortern(t *testing.T) {
	mockCodec := &mockTypedCodec[int32]{decodeVal: 42}
	attrs := []string{"first", "second", "third"}
	codec := NewMultiCodec(attrs, mockCodec, false)
	buf := &bytes.Buffer{}

	values := map[string]int32{
		"first":  1,
		"second": 2,
		"third":  3,
	}

	n, err := codec.Encode(values, buf)
	assert.NoError(t, err)
	assert.Equal(t, 3*1, n) // 3 elements * 1 byte (from mock)
	assert.True(t, mockCodec.encodeCalled)

	result, err := codec.Decode(buf)
	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, int32(42), result["first"])
	assert.Equal(t, int32(42), result["second"])
	assert.Equal(t, int32(42), result["third"])
}

func TestMultiCodec_WithBoolshortern_NotEmpty(t *testing.T) {
	mockCodec := &mockTypedCodec[int32]{decodeVal: 42}
	attrs := []string{"x", "y"}
	codec := NewMultiCodec(attrs, mockCodec, true)
	buf := &bytes.Buffer{}

	values := map[string]int32{
		"x": 10,
		"y": 20,
	}

	n, err := codec.Encode(values, buf)
	assert.NoError(t, err)
	assert.Equal(t, 1+2*1, n) // empty flag(1) + 2 elements

	// First byte should be 0 (not empty)
	assert.Equal(t, byte(0x00), buf.Bytes()[0])

	result, err := codec.Decode(buf)
	assert.NoError(t, err)
	assert.Equal(t, int32(42), result["x"])
	assert.Equal(t, int32(42), result["y"])
}

func TestMultiCodec_WithBoolshortern_Empty(t *testing.T) {
	mockCodec := &mockTypedCodec[int32]{}
	attrs := []string{"x", "y"}
	codec := NewMultiCodec(attrs, mockCodec, true)
	buf := &bytes.Buffer{}

	values := map[string]int32{}

	n, err := codec.Encode(values, buf)
	assert.NoError(t, err)
	assert.Equal(t, 1, n) // only empty flag

	// First byte should be 1 (empty)
	assert.Equal(t, byte(0x01), buf.Bytes()[0])

	result, err := codec.Decode(buf)
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestMultiCodec_MissingAttribute(t *testing.T) {
	mockCodec := &mockTypedCodec[int32]{}
	attrs := []string{"a", "b", "c"}
	codec := NewMultiCodec(attrs, mockCodec, false)
	buf := &bytes.Buffer{}

	values := map[string]int32{
		"a": 1,
		"b": 2,
		// "c" is missing
	}

	_, err := codec.Encode(values, buf)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing attribute")
}

func TestMultiCodec_ErrorHandling(t *testing.T) {
	t.Run("decode with failing codec", func(t *testing.T) {
		mockCodec := &mockTypedCodec[int32]{decodeErr: assert.AnError}
		attrs := []string{"test"}
		codec := NewMultiCodec(attrs, mockCodec, false)
		buf := &bytes.Buffer{}

		// Write some data that will be read by the mock
		buf.WriteByte(0xAA)

		_, err := codec.Decode(buf)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decode")
	})

	t.Run("encode with failing codec", func(t *testing.T) {
		mockCodec := &mockTypedCodec[int32]{encodeErr: assert.AnError}
		attrs := []string{"test"}
		codec := NewMultiCodec(attrs, mockCodec, false)
		buf := &bytes.Buffer{}

		values := map[string]int32{"test": 42}

		_, err := codec.Encode(values, buf)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to encode")
	})
}

func TestMultiCodec_OrderPreservation(t *testing.T) {
	// Test that attributes are encoded/decoded in the given order
	intCodec := &primitive.IntCodec{}
	attrs := []string{"z", "a", "m"} // deliberately unsorted
	codec := NewMultiCodec(attrs, intCodec, false)
	buf := &bytes.Buffer{}

	values := map[string]int32{
		"z": 3,
		"a": 1,
		"m": 2,
	}

	// Encode
	_, err := codec.Encode(values, buf)
	assert.NoError(t, err)

	// Buffer should contain values in order: z, a, m
	bytes := buf.Bytes()
	assert.Len(t, bytes, 3*4) // 3 ints * 4 bytes

	// Decode into new map
	result, err := codec.Decode(buf)
	assert.NoError(t, err)

	// Values should be mapped to correct attributes
	assert.Equal(t, int32(3), result["z"])
	assert.Equal(t, int32(1), result["a"])
	assert.Equal(t, int32(2), result["m"])
}

func TestMultiCodec_WithRealTypes(t *testing.T) {
	t.Run("float multi", func(t *testing.T) {
		floatCodec := &primitive.FloatCodec{}
		attrs := []string{"x", "y", "z"}
		codec := NewMultiCodec(attrs, floatCodec, true)
		buf := &bytes.Buffer{}

		values := map[string]float32{
			"x": 1.5,
			"y": 2.5,
			"z": 3.5,
		}

		n, err := codec.Encode(values, buf)
		assert.NoError(t, err)
		assert.Equal(t, 1+3*4, n)

		result, err := codec.Decode(buf)
		assert.NoError(t, err)
		assert.InDelta(t, 1.5, result["x"], 0.0001)
		assert.InDelta(t, 2.5, result["y"], 0.0001)
		assert.InDelta(t, 3.5, result["z"], 0.0001)
	})

	t.Run("bool multi", func(t *testing.T) {
		boolCodec := &primitive.BoolCodec{}
		attrs := []string{"flag1", "flag2"}
		codec := NewMultiCodec(attrs, boolCodec, false)
		buf := &bytes.Buffer{}

		values := map[string]bool{
			"flag1": true,
			"flag2": false,
		}

		n, err := codec.Encode(values, buf)
		assert.NoError(t, err)
		assert.Equal(t, 2, n) // 2 bools

		result, err := codec.Decode(buf)
		assert.NoError(t, err)
		assert.Equal(t, true, result["flag1"])
		assert.Equal(t, false, result["flag2"])
	})
}
