package packets

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/the-new-day/probogo/internal/codec"
	"github.com/the-new-day/probogo/internal/codec/complex"
	"github.com/the-new-day/probogo/internal/codec/primitive"
	"github.com/the-new-day/probogo/internal/protection"
)

// --- Mocks and helpers ---

// mockProtection implements the Protection interface for testing
type mockProtection struct {
	encryptCalled bool
	decryptCalled bool
	encryptResult []byte
	decryptResult []byte
}

func (m *mockProtection) Encrypt(data []byte) []byte {
	m.encryptCalled = true
	if m.encryptResult != nil {
		return m.encryptResult
	}
	// Default: return a copy (simulate no change)
	result := make([]byte, len(data))
	copy(result, data)
	return result
}

func (m *mockProtection) Decrypt(data []byte) []byte {
	m.decryptCalled = true
	if m.decryptResult != nil {
		return m.decryptResult
	}
	result := make([]byte, len(data))
	copy(result, data)
	return result
}

// simpleCodec implements codec.Codec for testing
type simpleCodec struct {
	encodeValue  any
	encodeErr    error
	decodeValue  any
	decodeErr    error
	encodeCalled bool
	decodeCalled bool
}

func (c *simpleCodec) Encode(value any, buf *bytes.Buffer) (int, error) {
	c.encodeCalled = true
	c.encodeValue = value
	if c.encodeErr != nil {
		return 0, c.encodeErr
	}
	// Just write a marker byte
	return buf.Write([]byte{0xAA})
}

func (c *simpleCodec) Decode(buf *bytes.Buffer) (any, error) {
	c.decodeCalled = true
	if c.decodeErr != nil {
		return nil, c.decodeErr
	}
	// Consume a byte
	buf.ReadByte()
	return c.decodeValue, nil
}

// --- Tests ---

func TestNewBasePacket(t *testing.T) {
	codecs := []codec.Codec{codec.Wrap(&primitive.IntCodec{}), codec.Wrap(&primitive.BoolCodec{})}
	attrs := []string{"field1", "field2"}

	packet := NewBasePacket(123, codecs, attrs, true)

	assert.Equal(t, int32(123), packet.id)
	assert.Equal(t, codecs, packet.codecs)
	assert.Equal(t, attrs, packet.attributes)
	assert.Equal(t, true, packet.shouldLog)
	assert.Empty(t, packet.objects)
	assert.Empty(t, packet.object)
}

func TestNewBasePacket_PanicOnMismatch(t *testing.T) {
	codecs := []codec.Codec{codec.Wrap(&primitive.IntCodec{})}
	attrs := []string{"field1", "field2"}

	assert.PanicsWithValue(t,
		"NewBasePacket: codecs and attributes length must be equal; codecs has 1 elements, attributes has 2 elements",
		func() {
			NewBasePacket(1, codecs, attrs, true)
		})
}

func TestBasePacket_Unwrap_Success(t *testing.T) {
	// Setup
	mockCodec1 := &simpleCodec{decodeValue: int32(42)}
	mockCodec2 := &simpleCodec{decodeValue: true}
	codecs := []codec.Codec{mockCodec1, mockCodec2}
	attrs := []string{"number", "flag"}
	packet := NewBasePacket(1, codecs, attrs, true)

	// Data that will be "decoded" (one byte per codec)
	data := bytes.NewBuffer([]byte{0x01, 0x02})

	// Execute
	result, err := packet.Unwrap(data)

	// Assert
	assert.NoError(t, err)
	assert.True(t, mockCodec1.decodeCalled)
	assert.True(t, mockCodec2.decodeCalled)
	assert.Equal(t, int32(42), packet.objects[0])
	assert.Equal(t, true, packet.objects[1])
	assert.Equal(t, map[string]any{"number": int32(42), "flag": true}, result)
	assert.Equal(t, result, packet.object)
}

func TestBasePacket_Unwrap_Error(t *testing.T) {
	// Setup
	mockCodec := &simpleCodec{decodeErr: assert.AnError}
	codecs := []codec.Codec{mockCodec}
	packet := NewBasePacket(1, codecs, []string{"field"}, true)

	data := bytes.NewBuffer([]byte{0x01})

	// Execute
	result, err := packet.Unwrap(data)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unwrap")
	assert.Nil(t, result)
	assert.Empty(t, packet.objects)
}

func TestBasePacket_Wrap_Success(t *testing.T) {
	// Setup
	mockCodec1 := &simpleCodec{}
	mockCodec2 := &simpleCodec{}
	codecs := []codec.Codec{mockCodec1, mockCodec2}
	packet := NewBasePacket(123, codecs, []string{"f1", "f2"}, true)
	packet.objects = []any{int32(42), true}

	mockProt := &mockProtection{}

	// Execute
	result, err := packet.Wrap(mockProt)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify codecs were called
	assert.True(t, mockCodec1.encodeCalled)
	assert.True(t, mockCodec2.encodeCalled)
	assert.Equal(t, int32(42), mockCodec1.encodeValue)
	assert.Equal(t, true, mockCodec2.encodeValue)

	// Verify protection was called
	assert.True(t, mockProt.encryptCalled)
	assert.False(t, mockProt.decryptCalled)

	// Check packet structure: [4 bytes length][4 bytes ID][encrypted payload]
	bytes := result.Bytes()
	assert.Len(t, bytes, 4+4+2) // header(8) + 2 bytes from mock codecs

	// Check length field (should be HeaderLength + payload length = 8+2 = 10)
	length := int32(bytes[0])<<24 | int32(bytes[1])<<16 | int32(bytes[2])<<8 | int32(bytes[3])
	assert.Equal(t, int32(10), length)

	// Check ID field
	id := int32(bytes[4])<<24 | int32(bytes[5])<<16 | int32(bytes[6])<<8 | int32(bytes[7])
	assert.Equal(t, int32(123), id)
}

func TestBasePacket_Wrap_WithXorProtection(t *testing.T) {
	encryptProt := protection.NewXorProtection(false)
	decryptProt := protection.NewXorProtection(true)
	decryptProt.Activate([]byte{0x12, 0x34})
	encryptProt.Activate([]byte{0x12, 0x34})

	intCodec := &primitive.IntCodec{}
	packet := NewBasePacket(42, []codec.Codec{codec.Wrap(intCodec)}, []string{"value"}, true)
	packet.objects = []any{int32(1000)}

	result, err := packet.Wrap(encryptProt)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	encryptedPayload := result.Bytes()[8:]
	decrypted := decryptProt.Decrypt(encryptedPayload)

	newPacket := NewBasePacket(42, []codec.Codec{codec.Wrap(intCodec)}, []string{"value"}, true)
	unwrapped, err := newPacket.Unwrap(bytes.NewBuffer(decrypted))
	assert.NoError(t, err)

	assert.Equal(t, int32(1000), unwrapped["value"])
}

func TestBasePacket_Wrap_ErrorOnEncode(t *testing.T) {
	// Setup codec that fails on Encode
	mockCodec := &simpleCodec{encodeErr: assert.AnError}
	packet := NewBasePacket(1, []codec.Codec{mockCodec}, []string{"field"}, true)
	packet.objects = []any{int32(42)}

	mockProt := &mockProtection{}

	// Execute
	result, err := packet.Wrap(mockProt)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to encode")
	assert.Nil(t, result)
	assert.False(t, mockProt.encryptCalled) // Should not reach encryption
}

func TestBasePacket_Wrap_PanicOnNilProtection(t *testing.T) {
	mockCodec := &simpleCodec{}
	packet := NewBasePacket(1, []codec.Codec{mockCodec}, []string{"field"}, true)
	assert.Panics(t, func() {
		packet.Wrap(nil)
	})
}

func TestBasePacket_Implement(t *testing.T) {
	codecs := []codec.Codec{
		codec.Wrap(&primitive.IntCodec{}),
		codec.Wrap(complex.NewStringCodec()),
		codec.Wrap(&primitive.BoolCodec{}),
	}
	packet := NewBasePacket(1, codecs, []string{"a", "b", "c"}, true)
	packet.objects = []any{1, "hello", true}

	result := packet.implement()

	expected := map[string]any{
		"a": 1,
		"b": "hello",
		"c": true,
	}
	assert.Equal(t, expected, result)
	assert.Equal(t, expected, packet.object)
}

func TestBasePacket_Implement_Empty(t *testing.T) {
	packet := NewBasePacket(1, nil, []string{}, true)
	packet.objects = []any{}

	result := packet.implement()

	assert.Empty(t, result)
	assert.Empty(t, packet.object)
}

func TestBasePacket_Unwrap_Wrap_Integration(t *testing.T) {
	// Full roundtrip test with real codecs
	intCodec := &primitive.IntCodec{}
	boolCodec := &primitive.BoolCodec{}
	stringCodec := &mockStringCodec{}

	codecs := []codec.Codec{codec.Wrap(intCodec), codec.Wrap(boolCodec), stringCodec}
	attrs := []string{"age", "active", "name"}

	packet := NewBasePacket(555, codecs, attrs, true)
	originalData := map[string]any{
		"age":    int32(25),
		"active": true,
		"name":   "Alice",
	}

	// Manually set objects (like after deimplement)
	packet.objects = []any{int32(25), true, "Alice"}

	// Wrap
	mockProt := &mockProtection{}
	wrapped, err := packet.Wrap(mockProt)
	assert.NoError(t, err)

	// Create new packet for unwrapping
	newPacket := NewBasePacket(555, codecs, attrs, true)

	// Unwrap (skip header)
	payload := wrapped.Bytes()[8:] // after header
	decrypted := mockProt.Decrypt(payload)
	unwrapped, err := newPacket.Unwrap(bytes.NewBuffer(decrypted))

	assert.NoError(t, err)
	assert.Equal(t, originalData, unwrapped)
}

// Helper mock for string codec
type mockStringCodec struct{}

func (m *mockStringCodec) Encode(value any, buf *bytes.Buffer) (int, error) {
	str, ok := value.(string)
	if !ok {
		return 0, assert.AnError
	}
	return buf.Write([]byte(str))
}

func (m *mockStringCodec) Decode(buf *bytes.Buffer) (any, error) {
	// Simple decode: read all remaining bytes as string
	return buf.String(), nil
}
