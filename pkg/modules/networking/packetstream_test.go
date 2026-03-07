package networking

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/the-new-day/protanki-networking/pkg/modules/protection"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

var emptyOnActivateProtection = func(keys []byte) {}

// --- Mocks ---

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
	decryptFunc    func([]byte) []byte
	encryptFunc    func([]byte) []byte
	activateCalled bool
	activateKeys   []byte
}

func (m *mockProtection) Activate(keys []byte) {
	m.activateCalled = true
	m.activateKeys = keys
}

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
	// Default mock behavior: prepend "enc_" just to verify it was called
	return append([]byte("enc_"), data...)
}

type mockPacket struct {
	packets.BasePacket
	id                int32
	wrapErr           error
	unwrapErr         error
	unwrapCalled      bool
	setCompressCalled bool
}

func (m *mockPacket) ID() int32 { return m.id }

func (m *mockPacket) Attr(attribute string) any {
	if attribute == "keys" {
		return []byte{0x13, 0x37} // Mock keys for ActivateProtection test
	}
	return 0
}

func (m *mockPacket) Set(name string, value any) {}

func (m *mockPacket) Wrap(protection protection.Protection) (*bytes.Buffer, error) {
	if m.wrapErr != nil {
		return nil, m.wrapErr
	}
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

func (m *mockPacket) SetCompress(shouldCompress bool) {
	m.setCompressCalled = shouldCompress
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

// --- Tests ---

func TestPacketStream_SendVariants(t *testing.T) {

	tests := []struct {
		name        string
		sendFunc    func(*PacketStream) error
		conn        *mockConnection
		expectWrite int
		expectErr   error
	}{
		{
			name: "Send packet",
			conn: &mockConnection{},
			sendFunc: func(ps *PacketStream) error {
				return ps.Send(&mockPacket{id: 1001})
			},
			expectWrite: 1,
		},
		{
			name: "Send write error",
			conn: &mockConnection{writeErr: errors.New("write failed")},
			sendFunc: func(ps *PacketStream) error {
				return ps.Send(&mockPacket{id: 1001})
			},
			expectErr: errors.New("write failed"),
		},
		{
			name: "SendRaw",
			conn: &mockConnection{},
			sendFunc: func(ps *PacketStream) error {
				return ps.SendRaw([]byte{1, 2, 3})
			},
			expectWrite: 1,
		},
		{
			name: "SendRawEncrypted",
			conn: &mockConnection{},
			sendFunc: func(ps *PacketStream) error {
				return ps.SendRawEncrypted([]byte("secret"))
			},
			expectWrite: 1,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			ps := NewPacketStream(
				tt.conn,
				&mockProtection{},
				packets.NewPacketRegistry(),
				emptyOnActivateProtection,
			)

			err := tt.sendFunc(ps)

			if tt.expectErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Len(t, tt.conn.writeData, tt.expectWrite)
		})
	}
}

func TestPacketStream_InboundScenarios(t *testing.T) {

	tests := []struct {
		name          string
		packetID      int32
		data          []byte
		compressed    bool
		registrySetup func(*packets.PacketRegistry)
		expectError   bool
		expectPacket  bool
		expectUnknown bool
	}{
		{
			name:         "valid packet",
			packetID:     1001,
			data:         []byte{1, 2, 3, 4},
			expectPacket: true,
			registrySetup: func(r *packets.PacketRegistry) {
				r.Register(1001, "Test", func() packets.Packet {
					return &mockPacket{id: 1001}
				})
			},
		},
		{
			name:          "unknown packet",
			packetID:      9999,
			data:          []byte{1, 2, 3},
			expectPacket:  true,
			expectUnknown: true,
		},
		{
			name:        "decompression error",
			packetID:    1001,
			data:        []byte{1, 2, 3},
			compressed:  true,
			expectError: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			packetLen := int32(packets.HeaderLength + len(tt.data))

			call := 0
			conn := &mockConnection{
				onRead: func(n int) ([]byte, error) {
					call++
					switch call {
					case 1:
						return createTestHeader(packetLen, tt.packetID, tt.compressed)[0:4], nil
					case 2:
						return createTestHeader(packetLen, tt.packetID, tt.compressed)[4:8], nil
					case 3:
						return tt.data, nil
					default:
						return nil, io.EOF
					}
				},
			}

			reg := packets.NewPacketRegistry()
			if tt.registrySetup != nil {
				tt.registrySetup(reg)
			}

			ps := NewPacketStream(conn, &mockProtection{}, reg, emptyOnActivateProtection)

			ch := ps.Inbound(context.Background())
			res := <-ch

			if tt.expectError {
				assert.Error(t, res.Err)
				return
			}

			assert.NoError(t, res.Err)

			if tt.expectPacket {
				assert.NotNil(t, res.Packet)
			}

			if tt.expectUnknown {
				_, ok := res.Packet.(*packets.UnknownPacket)
				assert.True(t, ok)
			}

			assert.Equal(t, tt.packetID, res.ID)
			assert.Equal(t, packetLen, res.Length)
			assert.Equal(t, tt.compressed, res.WasCompressed)
		})
	}
}

func TestPacketStream_ActivateProtectionHook(t *testing.T) {
	tests := []struct {
		name       string
		packetID   int32
		expectHook bool
	}{
		{
			name:       "activate protection packet",
			packetID:   packets.ActivateProtectionID,
			expectHook: true,
		},
		{
			name:       "normal packet",
			packetID:   1001,
			expectHook: false,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			packetLen := int32(packets.HeaderLength + 1)

			reg := packets.NewPacketRegistry()
			reg.Register(tt.packetID, "Test", func() packets.Packet {
				return &mockPacket{id: tt.packetID}
			})

			call := 0
			conn := &mockConnection{
				onRead: func(n int) ([]byte, error) {
					call++
					switch call {
					case 1:
						return createTestHeader(packetLen, tt.packetID, false)[0:4], nil
					case 2:
						return createTestHeader(packetLen, tt.packetID, false)[4:8], nil
					case 3:
						return []byte{0}, nil
					default:
						return nil, io.EOF
					}
				},
			}

			hookCalled := false
			ps := NewPacketStream(
				conn,
				&mockProtection{},
				reg,
				func(keys []byte) { hookCalled = true },
			)

			res := <-ps.Inbound(context.Background())

			assert.NoError(t, res.Err)
			assert.Equal(t, tt.expectHook, hookCalled)
		})
	}
}
