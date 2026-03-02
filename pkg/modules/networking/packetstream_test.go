package networking

import (
	"bytes"
	"compress/flate"
	"context"
	"encoding/binary"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/the-new-day/probogo/pkg/modules/protection"
	"github.com/the-new-day/probogo/pkg/packets"
)

type mockConnection struct {
	readData  [][]byte
	readIndex int
	readErr   error
	writeData [][]byte
	writeErr  error
	onRead    func(n int) ([]byte, error)
}

func (m *mockConnection) Read(n int) ([]byte, error) {
	if m.onRead != nil {
		return m.onRead(n)
	}
	if m.readErr != nil {
		return nil, m.readErr
	}
	if m.readIndex >= len(m.readData) {
		return nil, io.EOF
	}
	data := m.readData[m.readIndex]
	m.readIndex++
	if len(data) != n {
		return nil, io.ErrUnexpectedEOF
	}
	return data, nil
}

func (m *mockConnection) Write(data []byte) (int, error) {
	if m.writeErr != nil {
		return 0, m.writeErr
	}
	m.writeData = append(m.writeData, data)
	return len(data), nil
}

func (m *mockConnection) Close() error                       { return nil }
func (m *mockConnection) SetDeadline(t time.Time) error      { return nil }
func (m *mockConnection) SetReadDeadline(t time.Time) error  { return nil }
func (m *mockConnection) SetWriteDeadline(t time.Time) error { return nil }

type mockProtection struct {
	decryptFunc func([]byte) []byte
	encryptFunc func([]byte) []byte
}

func (m *mockProtection) Activate(keys []byte) {}

func (m *mockProtection) Decrypt(data []byte) []byte {
	if m.decryptFunc != nil {
		return m.decryptFunc(data)
	}
	return data
}

func (m *mockProtection) Encrypt(data []byte) []byte {
	if m.encryptFunc != nil {
		return m.encryptFunc(data)
	}
	return data
}

type mockPacket struct {
	id           int32
	wrapErr      error
	unwrapErr    error
	unwrapCalled bool
}

func (m *mockPacket) ID() int32 { return m.id }

func (m *mockPacket) Get(attribute string) any { return 0 }

func (m *mockPacket) Wrap(protection protection.Protection) (*bytes.Buffer, error) {
	if m.wrapErr != nil {
		return nil, m.wrapErr
	}
	// Simple header: length(8+data) + id + empty data
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, int32(8))
	binary.Write(buf, binary.BigEndian, m.id)
	return buf, nil
}

func (m *mockPacket) Unwrap(buf *bytes.Buffer) (map[string]any, error) {
	m.unwrapCalled = true
	if m.unwrapErr != nil {
		return nil, m.unwrapErr
	}
	return nil, nil
}

func createTestHeader(packetLen int32, packetID int32, compressed bool) []byte {
	header := make([]byte, 8)
	lenPart := uint32(packetLen) & 0x00FFFFFF
	if compressed {
		lenPart |= 0x40 << 24
	}
	binary.BigEndian.PutUint32(header[0:4], lenPart)
	binary.BigEndian.PutUint32(header[4:8], uint32(packetID))
	return header
}

func TestSend_Success(t *testing.T) {
	mockConn := &mockConnection{}
	mockProt := &mockProtection{}
	reg := packets.NewPacketRegistry()

	stream := NewPacketStream(mockConn, mockProt, reg)
	err := stream.Send(&mockPacket{id: 1001})

	assert.NoError(t, err)
	assert.Len(t, mockConn.writeData, 1)
}

func TestSend_WriteError(t *testing.T) {
	expectedErr := errors.New("write failed")
	mockConn := &mockConnection{writeErr: expectedErr}
	mockProt := &mockProtection{}
	reg := packets.NewPacketRegistry()

	stream := NewPacketStream(mockConn, mockProt, reg)
	err := stream.Send(&mockPacket{id: 1001})

	assert.ErrorIs(t, err, expectedErr)
}

func TestSend_WrapError(t *testing.T) {
	mockConn := &mockConnection{}
	mockProt := &mockProtection{}
	reg := packets.NewPacketRegistry()

	stream := NewPacketStream(mockConn, mockProt, reg)
	err := stream.Send(&mockPacket{id: 1001, wrapErr: errors.New("wrap failed")})

	assert.Error(t, err)
}

func TestPackets_ContextCancel(t *testing.T) {
	mockConn := &mockConnection{}
	mockProt := &mockProtection{}
	reg := packets.NewPacketRegistry()

	stream := NewPacketStream(mockConn, mockProt, reg)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	ch := stream.Packets(ctx)
	_, ok := <-ch
	assert.False(t, ok)
}

func TestPackets_ReadHeaderError(t *testing.T) {
	expectedErr := errors.New("read error")
	callCount := 0
	mockConn := &mockConnection{
		onRead: func(n int) ([]byte, error) {
			callCount++
			if callCount == 1 {
				return nil, expectedErr
			}
			return nil, nil
		},
	}
	mockProt := &mockProtection{}
	reg := packets.NewPacketRegistry()

	stream := NewPacketStream(mockConn, mockProt, reg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := stream.Packets(ctx)
	res := <-ch

	assert.ErrorIs(t, res.Err, expectedErr)
	assert.Nil(t, res.Packet)
}

func TestPackets_Success(t *testing.T) {
	packetID := int32(1001)
	packetData := []byte{1, 2, 3, 4}
	packetLen := int32(packetHeaderLength + len(packetData))

	reg := packets.NewPacketRegistry()
	reg.Register(packetID, "TestPacket", func() packets.Packet {
		return &mockPacket{id: packetID}
	})

	callCount := 0
	mockConn := &mockConnection{
		onRead: func(n int) ([]byte, error) {
			callCount++
			switch callCount {
			case 1:
				return createTestHeader(packetLen, packetID, false)[0:4], nil
			case 2:
				return createTestHeader(packetLen, packetID, false)[4:8], nil
			case 3:
				return packetData, nil
			default:
				return nil, io.EOF
			}
		},
	}

	mockProt := &mockProtection{}
	stream := NewPacketStream(mockConn, mockProt, reg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := stream.Packets(ctx)
	res := <-ch

	assert.NoError(t, res.Err)
	assert.NotNil(t, res.Packet)
	assert.Equal(t, packetID, res.Packet.ID())
}

func TestPackets_CompressedPacket(t *testing.T) {
	packetID := int32(1001)
	originalData := []byte("test data")
	var compressedData bytes.Buffer
	w, _ := flate.NewWriter(&compressedData, flate.DefaultCompression)
	w.Write(originalData)
	w.Close()
	packetLen := int32(packetHeaderLength + compressedData.Len())

	reg := packets.NewPacketRegistry()
	reg.Register(packetID, "TestPacket", func() packets.Packet {
		return &mockPacket{id: packetID}
	})

	callCount := 0
	mockConn := &mockConnection{
		onRead: func(n int) ([]byte, error) {
			callCount++
			switch callCount {
			case 1:
				return createTestHeader(packetLen, packetID, true)[0:4], nil
			case 2:
				return createTestHeader(packetLen, packetID, true)[4:8], nil
			case 3:
				return compressedData.Bytes(), nil
			default:
				return nil, io.EOF
			}
		},
	}

	mockProt := &mockProtection{}
	stream := NewPacketStream(mockConn, mockProt, reg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := stream.Packets(ctx)
	res := <-ch

	assert.NoError(t, res.Err)
	assert.NotNil(t, res.Packet)
	assert.Equal(t, packetID, res.Packet.ID())
}

func TestPackets_UnknownPacket(t *testing.T) {
	unknownID := int32(9999)
	packetData := []byte{1, 2, 3, 4}
	packetLen := int32(packetHeaderLength + len(packetData))

	reg := packets.NewPacketRegistry() // empty registry

	callCount := 0
	mockConn := &mockConnection{
		onRead: func(n int) ([]byte, error) {
			callCount++
			switch callCount {
			case 1:
				return createTestHeader(packetLen, unknownID, false)[0:4], nil
			case 2:
				return createTestHeader(packetLen, unknownID, false)[4:8], nil
			case 3:
				return packetData, nil
			default:
				return nil, io.EOF
			}
		},
	}

	mockProt := &mockProtection{}
	stream := NewPacketStream(mockConn, mockProt, reg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := stream.Packets(ctx)
	res := <-ch

	assert.NoError(t, res.Err)
	_, ok := res.Packet.(*packets.UnknownPacket)
	assert.True(t, ok, "expected UnknownPacket")
}

func TestPackets_UnwrapError(t *testing.T) {
	packetID := int32(1001)
	packetData := []byte{1, 2, 3, 4}
	packetLen := int32(packetHeaderLength + len(packetData))

	reg := packets.NewPacketRegistry()
	reg.Register(packetID, "TestPacket", func() packets.Packet {
		return &mockPacket{id: packetID, unwrapErr: errors.New("unwrap failed")}
	})

	callCount := 0
	mockConn := &mockConnection{
		onRead: func(n int) ([]byte, error) {
			callCount++
			switch callCount {
			case 1:
				return createTestHeader(packetLen, packetID, false)[0:4], nil
			case 2:
				return createTestHeader(packetLen, packetID, false)[4:8], nil
			case 3:
				return packetData, nil
			default:
				return nil, io.EOF
			}
		},
	}

	mockProt := &mockProtection{}
	stream := NewPacketStream(mockConn, mockProt, reg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := stream.Packets(ctx)
	res := <-ch

	assert.Error(t, res.Err)
	assert.Nil(t, res.Packet)
}

func TestPackets_DecompressionError(t *testing.T) {
	packetID := int32(1001)
	packetLen := int32(packetHeaderLength + 10)

	reg := packets.NewPacketRegistry()
	reg.Register(packetID, "TestPacket", func() packets.Packet {
		return &mockPacket{id: packetID}
	})

	callCount := 0
	mockConn := &mockConnection{
		onRead: func(n int) ([]byte, error) {
			callCount++
			switch callCount {
			case 1:
				return createTestHeader(packetLen, packetID, true)[0:4], nil
			case 2:
				return createTestHeader(packetLen, packetID, true)[4:8], nil
			case 3:
				// Return invalid compressed data
				return []byte{0xFF, 0xFF, 0xFF, 0xFF}, nil
			default:
				return nil, io.EOF
			}
		},
	}

	mockProt := &mockProtection{}
	stream := NewPacketStream(mockConn, mockProt, reg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := stream.Packets(ctx)
	res := <-ch

	assert.Error(t, res.Err)
	assert.Contains(t, res.Err.Error(), "decompression failed")
	assert.Nil(t, res.Packet)
}
