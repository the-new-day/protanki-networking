package packets

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/primitive"
	"github.com/the-new-day/protanki-networking/pkg/modules/protection"
)

// --- Mocks and helpers ---

// mockProtection implements the Protection interface for testing
type mockProtection struct {
	encryptCalled bool
	decryptCalled bool
	encryptResult []byte
	decryptResult []byte
}

func (m *mockProtection) Activate(keys []byte) {}

func (m *mockProtection) Encrypt(data []byte) []byte {
	m.encryptCalled = true
	if m.encryptResult != nil {
		return m.encryptResult
	}
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
	return buf.String(), nil
}

// --- Tests ---

func TestNewBasePacket(t *testing.T) {
	codecs := []codec.Codec{codec.Wrap(&primitive.IntCodec{}), codec.Wrap(&primitive.BoolCodec{})}
	attrs := []string{"field1", "field2"}

	packet := NewBasePacket(123, codecs, attrs)

	assert.Equal(t, int32(123), packet.ID())
	assert.Equal(t, codecs, packet.codecs)
	assert.Equal(t, attrs, packet.attributes)
	assert.NotNil(t, packet.objects)
	assert.NotNil(t, packet.object)
	assert.Empty(t, packet.objects)
	assert.Empty(t, packet.object)
	assert.Empty(t, packet.Data())
	assert.Equal(t, 0, packet.Len())
	assert.False(t, packet.shouldCompress)
}

func TestNewBasePacket_PanicOnMismatch(t *testing.T) {
	codecs := []codec.Codec{codec.Wrap(&primitive.IntCodec{})}
	attrs := []string{"field1", "field2"}

	assert.PanicsWithValue(t,
		"NewBasePacket: codecs and attributes length must be equal; codecs has 1 elements, attributes has 2 elements",
		func() {
			NewBasePacket(1, codecs, attrs)
		})
}

func TestBasePacket_Unwrap_Success(t *testing.T) {
	mockCodec1 := &simpleCodec{decodeValue: int32(42)}
	mockCodec2 := &simpleCodec{decodeValue: true}
	codecs := []codec.Codec{mockCodec1, mockCodec2}
	attrs := []string{"number", "flag"}
	packet := NewBasePacket(1, codecs, attrs)

	// Data that will be "decoded" (one byte per codec)
	inputBytes := []byte{0x01, 0x02}
	data := bytes.NewBuffer(inputBytes)

	result, err := packet.Unwrap(data)

	assert.NoError(t, err)
	assert.True(t, mockCodec1.decodeCalled)
	assert.True(t, mockCodec2.decodeCalled)

	// Check objects array
	assert.Equal(t, int32(42), packet.objects[0])
	assert.Equal(t, true, packet.objects[1])

	// Check returned map and internal object map
	expectedMap := map[string]any{"number": int32(42), "flag": true}
	assert.Equal(t, expectedMap, result)
	assert.Equal(t, expectedMap, packet.object)

	// Check new Data() and Len() methods
	assert.Equal(t, inputBytes, packet.Data())
	assert.Equal(t, len(inputBytes), packet.Len())
}

func TestBasePacket_Unwrap_Error(t *testing.T) {
	mockCodec := &simpleCodec{decodeErr: assert.AnError}
	packet := NewBasePacket(1, []codec.Codec{mockCodec}, []string{"field"})

	data := bytes.NewBuffer([]byte{0x01})
	result, err := packet.Unwrap(data)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unwrap")
	assert.Nil(t, result)
	assert.Empty(t, packet.objects)
	// rawData should not be fully populated on error in the middle of decoding
	assert.Empty(t, packet.Data())
}

func TestBasePacket_Wrap_Success(t *testing.T) {
	mockCodec1 := &simpleCodec{}
	mockCodec2 := &simpleCodec{}
	packet := NewBasePacket(123, []codec.Codec{mockCodec1, mockCodec2}, []string{"f1", "f2"})
	packet.objects = []any{int32(42), true}

	mockProt := &mockProtection{}
	result, err := packet.Wrap(mockProt)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.True(t, mockCodec1.encodeCalled)
	assert.True(t, mockCodec2.encodeCalled)

	bytesData := result.Bytes()
	assert.Len(t, bytesData, 8+2) // header(8) + 2 bytes from mock codecs

	// Check length field (should be HeaderLength + payload length = 8+2 = 10)
	length := int32(bytesData[0])<<24 | int32(bytesData[1])<<16 | int32(bytesData[2])<<8 | int32(bytesData[3])
	assert.Equal(t, int32(10), length)

	// Check ID field
	id := int32(bytesData[4])<<24 | int32(bytesData[5])<<16 | int32(bytesData[6])<<8 | int32(bytesData[7])
	assert.Equal(t, int32(123), id)
}

func TestBasePacket_Wrap_WithCompression(t *testing.T) {
	// Note: This test assumes the package-level Compress function works correctly and returns some byte slice.
	mockCodec := &simpleCodec{}
	packet := NewBasePacket(1, []codec.Codec{mockCodec}, []string{"val"})
	packet.objects = []any{int32(1)}

	packet.SetCompress(true)
	assert.True(t, packet.shouldCompress)

	mockProt := &mockProtection{}
	result, err := packet.Wrap(mockProt)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	bytesData := result.Bytes()
	length := int32(bytesData[0])<<24 | int32(bytesData[1])<<16 | int32(bytesData[2])<<8 | int32(bytesData[3])

	// Check if the compression bit (0x40000000) is set
	assert.True(t, (length&0x40000000) != 0, "Compression bit should be set in length header")
}

func TestBasePacket_Wrap_PanicOnMissingObjects(t *testing.T) {
	mockCodec := &simpleCodec{}
	packet := NewBasePacket(1, []codec.Codec{mockCodec}, []string{"field"})
	// We intentionally do NOT populate packet.objects to simulate bad manual creation

	mockProt := &mockProtection{}
	assert.Panics(t, func() {
		packet.Wrap(mockProt)
	}, "Wrap should panic due to index out of bounds when accessing bp.objects[i]")
}

func TestBasePacket_Attr_And_Set(t *testing.T) {
	packet := NewBasePacket(10, nil, nil) // Empty codecs and attrs for manual testing

	// Test Set
	packet.Set("customAttr", "customValue")
	packet.Set("number", 100)

	// Test Attr
	assert.Equal(t, "customValue", packet.Attr("customAttr"))
	assert.Equal(t, 100, packet.Attr("number"))
}

func TestBasePacket_Attr_PanicOnMissing(t *testing.T) {
	packet := NewBasePacket(99, nil, nil)

	assert.PanicsWithValue(t,
		"BasePacket.Get: attribute \"missing\" not found. ID: 99",
		func() {
			packet.Attr("missing")
		},
	)
}

func TestBasePacket_Wrap_WithXorProtection(t *testing.T) {
	encryptProt := protection.NewXorProtection(false)
	decryptProt := protection.NewXorProtection(true)
	decryptProt.Activate([]byte{0x12, 0x34})
	encryptProt.Activate([]byte{0x12, 0x34})

	intCodec := &primitive.IntCodec{}
	packet := NewBasePacket(42, []codec.Codec{codec.Wrap(intCodec)}, []string{"value"})
	packet.objects = []any{int32(1000)}

	result, err := packet.Wrap(encryptProt)
	assert.NoError(t, err)

	encryptedPayload := result.Bytes()[8:]
	decrypted := decryptProt.Decrypt(encryptedPayload)

	newPacket := NewBasePacket(42, []codec.Codec{codec.Wrap(intCodec)}, []string{"value"})
	unwrapped, err := newPacket.Unwrap(bytes.NewBuffer(decrypted))

	assert.NoError(t, err)
	assert.Equal(t, int32(1000), unwrapped["value"])
	assert.Equal(t, int32(1000), newPacket.Attr("value"))
}

func TestBasePacket_Wrap_ErrorOnEncode(t *testing.T) {
	mockCodec := &simpleCodec{encodeErr: assert.AnError}
	packet := NewBasePacket(1, []codec.Codec{mockCodec}, []string{"field"})
	packet.objects = []any{int32(42)}

	mockProt := &mockProtection{}
	result, err := packet.Wrap(mockProt)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to encode")
	assert.Nil(t, result)
	assert.False(t, mockProt.encryptCalled)
}

func TestBasePacket_Wrap_PanicOnNilProtection(t *testing.T) {
	mockCodec := &simpleCodec{}
	packet := NewBasePacket(1, []codec.Codec{mockCodec}, []string{"field"})

	assert.PanicsWithValue(t, "BasePacket.Wrap: nil protection is passed", func() {
		packet.Wrap(nil)
	})
}
