package primitive

import (
	"bytes"
	"io"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/the-new-day/protanki-networking/pkg/codec"
)

func TestBoolCodec(t *testing.T) {
	typed := &BoolCodec{}

	tests := []struct {
		name  string
		value bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}

			// Test TypedCodec
			n, err := typed.Encode(tt.value, buf)
			assert.NoError(t, err)
			assert.Equal(t, 1, n)
			assert.Equal(t, 1, buf.Len())

			result, err := typed.Decode(buf)
			assert.NoError(t, err)
			assert.Equal(t, tt.value, result)
			assert.Equal(t, 0, buf.Len())
		})
	}
}

func TestBoolCodec_AsCodec(t *testing.T) {
	typed := &BoolCodec{}
	wrapped := codec.Wrap(typed)
	buf := &bytes.Buffer{}

	// Test as Codec interface
	n, err := wrapped.Encode(true, buf)
	assert.NoError(t, err)
	assert.Equal(t, 1, n)

	result, err := wrapped.Decode(buf)
	assert.NoError(t, err)
	assert.Equal(t, true, result)
}

func TestBoolCodec_WrongTypePanic(t *testing.T) {
	typed := &BoolCodec{}
	wrapped := codec.Wrap(typed)
	buf := &bytes.Buffer{}

	assert.Panics(t, func() {
		wrapped.Encode("not a bool", buf)
	})
}

func TestByteCodec(t *testing.T) {
	typed := &ByteCodec{}

	tests := []struct {
		name  string
		value byte
	}{
		{"zero", 0},
		{"one", 1},
		{"max", 255},
		{"mid", 128},
		{"random", 0xAB},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}

			// Test TypedCodec
			n, err := typed.Encode(tt.value, buf)
			assert.NoError(t, err)
			assert.Equal(t, 1, n)
			assert.Equal(t, 1, buf.Len())

			result, err := typed.Decode(buf)
			assert.NoError(t, err)
			assert.Equal(t, tt.value, result)
			assert.Equal(t, 0, buf.Len())

			// Test wrapped
			wrapped := codec.Wrap(typed)
			buf2 := &bytes.Buffer{}
			n, err = wrapped.Encode(tt.value, buf2)
			assert.NoError(t, err)
			assert.Equal(t, 1, n)

			anyResult, err := wrapped.Decode(buf2)
			assert.NoError(t, err)
			assert.Equal(t, tt.value, anyResult)
		})
	}
}

func TestFloatCodec(t *testing.T) {
	typed := &FloatCodec{}
	epsilon := 0.0001

	tests := []struct {
		name  string
		value float32
	}{
		{"zero", 0.0},
		{"positive", 3.14},
		{"negative", -3.14},
		{"small", 1.5e-10},
		{"max", math.MaxFloat32},
		{"min positive", math.SmallestNonzeroFloat32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}

			// Test TypedCodec
			n, err := typed.Encode(tt.value, buf)
			assert.NoError(t, err)
			assert.Equal(t, 4, n)
			assert.Equal(t, 4, buf.Len())

			result, err := typed.Decode(buf)
			assert.NoError(t, err)
			assert.InDelta(t, tt.value, result, epsilon)
			assert.Equal(t, 0, buf.Len())
		})
	}
}

func TestFloatCodec_SpecialValues(t *testing.T) {
	typed := &FloatCodec{}
	buf := &bytes.Buffer{}

	tests := []struct {
		name  string
		value float32
		check func(t *testing.T, val float32)
	}{
		{
			name:  "NaN",
			value: float32(math.NaN()),
			check: func(t *testing.T, val float32) {
				assert.True(t, math.IsNaN(float64(val)))
			},
		},
		{
			name:  "+Inf",
			value: float32(math.Inf(1)),
			check: func(t *testing.T, val float32) {
				assert.True(t, math.IsInf(float64(val), 1))
			},
		},
		{
			name:  "-Inf",
			value: float32(math.Inf(-1)),
			check: func(t *testing.T, val float32) {
				assert.True(t, math.IsInf(float64(val), -1))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()

			_, err := typed.Encode(tt.value, buf)
			assert.NoError(t, err)

			result, err := typed.Decode(buf)
			assert.NoError(t, err)

			tt.check(t, result)
		})
	}
}

func TestFloatCodec_BigEndian(t *testing.T) {
	typed := &FloatCodec{}
	buf := &bytes.Buffer{}

	typed.Encode(1.0, buf)

	expected := []byte{0x3F, 0x80, 0x00, 0x00}
	assert.Equal(t, expected, buf.Bytes())
}

func TestIntCodec(t *testing.T) {
	typed := &IntCodec{}

	tests := []struct {
		name  string
		value int32
	}{
		{"zero", 0},
		{"positive", 42},
		{"negative", -42},
		{"max", 2147483647},
		{"min", -2147483648},
		{"hex", 0x12345678},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}

			n, err := typed.Encode(tt.value, buf)
			assert.NoError(t, err)
			assert.Equal(t, 4, n)

			result, err := typed.Decode(buf)
			assert.NoError(t, err)
			assert.Equal(t, tt.value, result)
			assert.Equal(t, 0, buf.Len())
		})
	}
}

func TestIntCodec_BigEndian(t *testing.T) {
	typed := &IntCodec{}
	buf := &bytes.Buffer{}

	typed.Encode(0x12345678, buf)

	expected := []byte{0x12, 0x34, 0x56, 0x78}
	assert.Equal(t, expected, buf.Bytes())
}

func TestShortCodec(t *testing.T) {
	typed := &ShortCodec{}

	tests := []struct {
		name  string
		value int16
	}{
		{"zero", 0},
		{"positive", 42},
		{"negative", -42},
		{"max", 32767},
		{"min", -32768},
		{"hex", 0x1234},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}

			n, err := typed.Encode(tt.value, buf)
			assert.NoError(t, err)
			assert.Equal(t, 2, n)

			result, err := typed.Decode(buf)
			assert.NoError(t, err)
			assert.Equal(t, tt.value, result)
			assert.Equal(t, 0, buf.Len())
		})
	}
}

func TestShortCodec_BigEndian(t *testing.T) {
	typed := &ShortCodec{}
	buf := &bytes.Buffer{}

	typed.Encode(0x1234, buf)

	expected := []byte{0x12, 0x34}
	assert.Equal(t, expected, buf.Bytes())
}

func TestLongCodec(t *testing.T) {
	typed := &LongCodec{}

	tests := []struct {
		name  string
		value int64
	}{
		{"zero", 0},
		{"positive", 42},
		{"negative", -42},
		{"max", 9223372036854775807},
		{"min", -9223372036854775808},
		{"hex", 0x123456789ABCDEF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}

			n, err := typed.Encode(tt.value, buf)
			assert.NoError(t, err)
			assert.Equal(t, 8, n)

			result, err := typed.Decode(buf)
			assert.NoError(t, err)
			assert.Equal(t, tt.value, result)
			assert.Equal(t, 0, buf.Len())
		})
	}
}

func TestLongCodec_BigEndian(t *testing.T) {
	typed := &LongCodec{}
	buf := &bytes.Buffer{}

	typed.Encode(0x123456789ABCDEF, buf)

	expected := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF}
	assert.Equal(t, expected, buf.Bytes())
}

func TestErrorCases(t *testing.T) {
	t.Run("empty reader", func(t *testing.T) {
		empty := &bytes.Buffer{}

		_, err := (&BoolCodec{}).Decode(empty)
		assert.ErrorIs(t, err, io.EOF)

		_, err = (&ByteCodec{}).Decode(empty)
		assert.ErrorIs(t, err, io.EOF)

		_, err = (&FloatCodec{}).Decode(empty)
		assert.ErrorIs(t, err, io.EOF)

		_, err = (&IntCodec{}).Decode(empty)
		assert.ErrorIs(t, err, io.EOF)

		_, err = (&ShortCodec{}).Decode(empty)
		assert.ErrorIs(t, err, io.EOF)

		_, err = (&LongCodec{}).Decode(empty)
		assert.ErrorIs(t, err, io.EOF)
	})

	t.Run("insufficient data", func(t *testing.T) {
		partial := bytes.NewBuffer([]byte{0x01, 0x02})

		_, err := (&IntCodec{}).Decode(partial)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected EOF")

		partial.Reset()
		_, err = (&ShortCodec{}).Decode(partial)
		assert.Error(t, err)

		partial.Reset()
		_, err = (&LongCodec{}).Decode(partial)
		assert.Error(t, err)
	})
}

func TestMultipleValues(t *testing.T) {
	buf := &bytes.Buffer{}

	boolCodec := &BoolCodec{}
	byteCodec := &ByteCodec{}
	intCodec := &IntCodec{}
	shortCodec := &ShortCodec{}
	longCodec := &LongCodec{}
	floatCodec := &FloatCodec{}

	// Write various values
	intCodec.Encode(42, buf)
	boolCodec.Encode(true, buf)
	shortCodec.Encode(1234, buf)
	floatCodec.Encode(3.14, buf)
	byteCodec.Encode(0xFF, buf)
	longCodec.Encode(999999, buf)

	// Read them back
	v1, _ := intCodec.Decode(buf)
	assert.Equal(t, int32(42), v1)

	v2, _ := boolCodec.Decode(buf)
	assert.Equal(t, true, v2)

	v3, _ := shortCodec.Decode(buf)
	assert.Equal(t, int16(1234), v3)

	v4, _ := floatCodec.Decode(buf)
	assert.InDelta(t, 3.14, v4, 0.0001)

	v5, _ := byteCodec.Decode(buf)
	assert.Equal(t, byte(0xFF), v5)

	v6, _ := longCodec.Decode(buf)
	assert.Equal(t, int64(999999), v6)

	assert.Equal(t, 0, buf.Len())
}

func TestInterfaceCompliance(t *testing.T) {
	// Verify all codecs implement TypedCodec
	var _ codec.TypedCodec[bool] = &BoolCodec{}
	var _ codec.TypedCodec[byte] = &ByteCodec{}
	var _ codec.TypedCodec[float32] = &FloatCodec{}
	var _ codec.TypedCodec[int32] = &IntCodec{}
	var _ codec.TypedCodec[int16] = &ShortCodec{}
	var _ codec.TypedCodec[int64] = &LongCodec{}

	// Verify wrappers implement Codec
	var _ codec.Codec = codec.Wrap(&BoolCodec{})
	var _ codec.Codec = codec.Wrap(&ByteCodec{})
	var _ codec.Codec = codec.Wrap(&FloatCodec{})
	var _ codec.Codec = codec.Wrap(&IntCodec{})
	var _ codec.Codec = codec.Wrap(&ShortCodec{})
	var _ codec.Codec = codec.Wrap(&LongCodec{})
}

func TestCodecWrapper_Array(t *testing.T) {
	codecs := []codec.Codec{
		codec.Wrap(&BoolCodec{}),
		codec.Wrap(&IntCodec{}),
	}

	assert.Len(t, codecs, 2)
}
