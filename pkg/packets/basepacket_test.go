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
	res := make([]byte, len(data))
	copy(res, data)
	return res
}

func (m *mockProtection) Decrypt(data []byte) []byte {
	m.decryptCalled = true
	if m.decryptResult != nil {
		return m.decryptResult
	}
	res := make([]byte, len(data))
	copy(res, data)
	return res
}

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
	return buf.Write([]byte{0xAA})
}

func (c *simpleCodec) Decode(buf *bytes.Buffer) (any, error) {
	c.decodeCalled = true
	if c.decodeErr != nil {
		return nil, c.decodeErr
	}
	_, _ = buf.ReadByte()
	return c.decodeValue, nil
}

type mockStringCodec struct{}

func (m *mockStringCodec) Encode(value any, buf *bytes.Buffer) (int, error) {
	s, ok := value.(string)
	if !ok {
		return 0, assert.AnError
	}
	return buf.Write([]byte(s))
}

func (m *mockStringCodec) Decode(buf *bytes.Buffer) (any, error) {
	return buf.String(), nil
}

// --- Tests ---

func TestNewBasePacket_ConstructorsAndPanics(t *testing.T) {
	t.Run("NewBasePacket - happy path", func(t *testing.T) {
		codecs := []codec.Codec{codec.Wrap(&primitive.IntCodec{}), codec.Wrap(&primitive.BoolCodec{})}
		attrs := []string{"field1", "field2"}

		packet := NewBasePacket(123, codecs, attrs)

		assert.Equal(t, int32(123), packet.ID())
		assert.Empty(t, packet.Data())
		assert.Equal(t, 0, packet.Len())
	})

	t.Run("NewBasePacket - mismatch panic", func(t *testing.T) {
		codecs := []codec.Codec{codec.Wrap(&primitive.IntCodec{})}
		attrs := []string{"field1", "field2"}

		assert.PanicsWithValue(t,
			"NewBasePacket: codecs and attributes length must be equal; codecs has 1 elements, attributes has 2 elements",
			func() { NewBasePacket(1, codecs, attrs) })
	})
}

func TestBasePacket_UnwrapAndUnwrapValues(t *testing.T) {
	t.Run("Unwrap - success populates objects and raw data", func(t *testing.T) {
		mock1 := &simpleCodec{decodeValue: int32(42)}
		mock2 := &simpleCodec{decodeValue: true}
		codecs := []codec.Codec{mock1, mock2}
		attrs := []string{"number", "flag"}

		packet := NewBasePacket(1, codecs, attrs)

		inputBytes := []byte{0x01, 0x02}
		data := bytes.NewBuffer(inputBytes)

		result, err := packet.Unwrap(data)
		assert.NoError(t, err)

		assert.True(t, mock1.decodeCalled)
		assert.True(t, mock2.decodeCalled)

		assert.Equal(t, int32(42), packet.objects[0])
		assert.Equal(t, true, packet.objects[1])

		expected := map[string]any{"number": int32(42), "flag": true}
		assert.Equal(t, expected, result)
		assert.Equal(t, expected, packet.object)

		assert.Equal(t, inputBytes, packet.Data())
		assert.Equal(t, len(inputBytes), packet.Len())
	})

	t.Run("Unwrap - error in decode leaves objects empty and returns error", func(t *testing.T) {
		mockErr := &simpleCodec{decodeErr: assert.AnError}
		packet := NewBasePacket(1, []codec.Codec{mockErr}, []string{"field"})

		data := bytes.NewBuffer([]byte{0x01})
		result, err := packet.Unwrap(data)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unwrap")
		assert.Nil(t, result)
		assert.Empty(t, packet.Data())
	})

	t.Run("UnwrapValues - encodes provided values and delegates to Unwrap", func(t *testing.T) {
		intC := &primitive.IntCodec{}
		packet := NewBasePacket(2, []codec.Codec{codec.Wrap(intC), codec.Wrap(intC)}, []string{"a", "b"})

		res, err := packet.UnwrapValues(int32(10), int32(20))
		assert.NoError(t, err)
		// After UnwrapValues, packet.objects populated (via Unwrap)
		assert.Equal(t, int32(10), packet.Attr("a"))
		assert.Equal(t, int32(20), packet.Attr("b"))
		assert.Equal(t, int32(10), res["a"])
		assert.Equal(t, int32(20), res["b"])
	})
}

func TestBasePacket_WrapVariants(t *testing.T) {
	t.Run("Wrap - success encodes header and payload", func(t *testing.T) {
		c1 := &simpleCodec{}
		c2 := &simpleCodec{}
		packet := NewBasePacket(123, []codec.Codec{c1, c2}, []string{"f1", "f2"})
		packet.UnwrapValues(int32(42), true)

		mockProt := &mockProtection{}
		out, err := packet.Wrap(mockProt)

		assert.NoError(t, err)
		assert.NotNil(t, out)
		assert.True(t, c1.encodeCalled)
		assert.True(t, c2.encodeCalled)

		bytesData := out.Bytes()
		// header (8) + 2 payload bytes (simpleCodec writes 1 byte each)
		assert.Len(t, bytesData, HeaderLength+2)

		// Check length field (big endian int32)
		length := int32(bytesData[0])<<24 | int32(bytesData[1])<<16 | int32(bytesData[2])<<8 | int32(bytesData[3])
		assert.Equal(t, int32(HeaderLength+2), length)

		id := int32(bytesData[4])<<24 | int32(bytesData[5])<<16 | int32(bytesData[6])<<8 | int32(bytesData[7])
		assert.Equal(t, int32(123), id)
	})

	t.Run("Wrap - set compress bit when compression enabled", func(t *testing.T) {
		c := &simpleCodec{}
		packet := NewBasePacket(1, []codec.Codec{c}, []string{"val"})
		packet.UnwrapValues(int32(1))

		packet.SetCompress(true)
		assert.True(t, packet.shouldCompress)

		mockProt := &mockProtection{}
		out, err := packet.Wrap(mockProt)
		assert.NoError(t, err)
		assert.NotNil(t, out)

		bytesData := out.Bytes()
		length := int32(bytesData[0])<<24 | int32(bytesData[1])<<16 | int32(bytesData[2])<<8 | int32(bytesData[3])

		// compression bit (0x40000000) must be set
		assert.True(t, (length&0x40000000) != 0, "Compression bit should be set in length header")
	})

	t.Run("Wrap - panic on nil protection", func(t *testing.T) {
		c := &simpleCodec{}
		packet := NewBasePacket(1, []codec.Codec{c}, []string{"field"})
		assert.PanicsWithValue(t, "BasePacket.Wrap: nil protection is passed", func() {
			_, _ = packet.Wrap(nil)
		})
	})

	t.Run("Wrap - panic on missing objects (index out of range)", func(t *testing.T) {
		c := &simpleCodec{}
		packet := NewBasePacket(1, []codec.Codec{c}, []string{"field"})
		// packet.objects intentionally not populated -> expected panic when accessing objects[i]
		assert.Panics(t, func() {
			_, _ = packet.Wrap(&mockProtection{})
		})
	})

	t.Run("Wrap - error when codec.Encode fails", func(t *testing.T) {
		c := &simpleCodec{encodeErr: assert.AnError}
		packet := NewBasePacket(1, []codec.Codec{c}, []string{"field"})
		packet.objects = []any{int32(42)}

		mockProt := &mockProtection{}
		out, err := packet.Wrap(mockProt)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to encode")
		assert.Nil(t, out)
		assert.False(t, mockProt.encryptCalled)
	})
}

func TestBasePacket_AttrAndSetBehavior(t *testing.T) {
	t.Run("Set then Attr returns values and Attr panics on missing", func(t *testing.T) {
		packet := NewBasePacket(10, nil, nil)

		packet.Set("customAttr", "customValue")
		packet.Set("number", 100)

		assert.Equal(t, "customValue", packet.Attr("customAttr"))
		assert.Equal(t, 100, packet.Attr("number"))

		assert.PanicsWithValue(t,
			"BasePacket.Get: attribute \"missing\" not found. ID: 10",
			func() { packet.Attr("missing") },
		)
	})
}

func TestBasePacket_Wrap_WithXorProtectionRoundtrip(t *testing.T) {
	t.Run("xor protection roundtrip encodes and decodes back", func(t *testing.T) {
		encryptProt := protection.NewXorProtection(false)
		decryptProt := protection.NewXorProtection(true)
		keys := []byte{0x12, 0x34}
		decryptProt.Activate(keys)
		encryptProt.Activate(keys)

		intCodec := &primitive.IntCodec{}
		packet := NewBasePacket(42, []codec.Codec{codec.Wrap(intCodec)}, []string{"value"})
		packet.UnwrapValues(int32(1000))

		out, err := packet.Wrap(encryptProt)
		assert.NoError(t, err)

		encryptedPayload := out.Bytes()[HeaderLength:]
		decrypted := decryptProt.Decrypt(encryptedPayload)

		newPacket := NewBasePacket(42, []codec.Codec{codec.Wrap(intCodec)}, []string{"value"})
		unwrapped, err := newPacket.Unwrap(bytes.NewBuffer(decrypted))
		assert.NoError(t, err)
		assert.Equal(t, int32(1000), unwrapped["value"])
		assert.Equal(t, int32(1000), newPacket.Attr("value"))
	})
}
