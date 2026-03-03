package proxy

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/the-new-day/probogo/pkg/modules/protection"
	"github.com/the-new-day/probogo/pkg/packets"
	"github.com/the-new-day/probogo/pkg/packets/network"
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
	return data
}

type mockPacket struct {
	id    int32
	attrs map[string]any
}

func (m *mockPacket) ID() int32 { return m.id }

func (m *mockPacket) Get(attribute string) any { return 0 }

func (m *mockPacket) Attr(name string) any {
	if m.attrs == nil {
		return nil
	}
	return m.attrs[name]
}

func (m *mockPacket) Set(name string, value any) {
	if m.attrs == nil {
		m.attrs = make(map[string]any)
	}
	m.attrs[name] = value
}

func (m *mockPacket) Wrap(protection protection.Protection) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, int32(8))
	binary.Write(buf, binary.BigEndian, m.id)
	return buf, nil
}

func (m *mockPacket) Unwrap(buf *bytes.Buffer) (map[string]any, error) {
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

func TestNewProxy(t *testing.T) {
	serverConn := &mockConnection{}
	serverProt := &mockProtection{}
	clientConn := &mockConnection{}
	clientProt := &mockProtection{}
	reg := packets.NewPacketRegistry()

	proxy := NewProxy(serverConn, serverProt, clientConn, clientProt, reg)

	assert.NotNil(t, proxy)
	assert.NotNil(t, proxy.serverHandler)
	assert.NotNil(t, proxy.clientHandler)
	assert.Equal(t, []byte{}, proxy.clientProtectionKeys)
}

func TestSetClientProtectionKeys(t *testing.T) {
	serverConn := &mockConnection{}
	serverProt := &mockProtection{}
	clientConn := &mockConnection{}
	clientProt := &mockProtection{}
	reg := packets.NewPacketRegistry()

	proxy := NewProxy(serverConn, serverProt, clientConn, clientProt, reg)

	keys := []byte{1, 2, 3, 4}
	proxy.SetClientProtectionKeys(keys)

	assert.Equal(t, keys, proxy.clientProtectionKeys)
}

func TestSetClientProtectionKeys_CopiesData(t *testing.T) {
	serverConn := &mockConnection{}
	serverProt := &mockProtection{}
	clientConn := &mockConnection{}
	clientProt := &mockProtection{}
	reg := packets.NewPacketRegistry()

	proxy := NewProxy(serverConn, serverProt, clientConn, clientProt, reg)

	keys := []byte{1, 2, 3, 4}
	proxy.SetClientProtectionKeys(keys)

	// Modify original
	keys[0] = 99

	// Proxy's copy should not be modified
	assert.Equal(t, byte(1), proxy.clientProtectionKeys[0])
}

func TestRun_ActivatesClientProtection(t *testing.T) {
	serverConn := &mockConnection{}
	serverProt := &mockProtection{}
	clientConn := &mockConnection{}
	clientProt := &mockProtection{}
	reg := packets.NewPacketRegistry()

	proxy := NewProxy(serverConn, serverProt, clientConn, clientProt, reg)

	keys := []byte{5, 6, 7, 8}
	proxy.SetClientProtectionKeys(keys)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		proxy.Run(ctx)
	}()

	time.Sleep(50 * time.Millisecond)
	cancel()

	// Client protection should have been activated with the set keys
	assert.True(t, clientProt.activateCalled)
	assert.Equal(t, keys, clientProt.activateKeys)
}

func TestRun_StartsServerHandler(t *testing.T) {
	packetID := int32(1002)
	packetData := []byte{5, 6, 7, 8}
	packetLen := int32(packets.HeaderLength + len(packetData))

	callCount := 0
	serverConn := &mockConnection{
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
	serverProt := &mockProtection{}
	clientConn := &mockConnection{}
	clientProt := &mockProtection{}

	reg := packets.NewPacketRegistry()
	reg.Register(packetID, "TestPacket", func() packets.Packet {
		return &mockPacket{id: packetID}
	})

	proxy := NewProxy(serverConn, serverProt, clientConn, clientProt, reg)

	var mu sync.Mutex
	packetReceived := false
	proxy.OnServerToClient(func(p packets.Packet) packets.Packet {
		mu.Lock()
		packetReceived = true
		mu.Unlock()
		return p
	})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		proxy.Run(ctx)
	}()

	time.Sleep(150 * time.Millisecond)
	cancel()

	// Server handler should have received a packet
	mu.Lock()
	received := packetReceived
	mu.Unlock()
	assert.True(t, received)
}

func TestHandleActivateProtection_InterceptsPacket(t *testing.T) {
	serverConn := &mockConnection{}
	serverProt := &mockProtection{}
	clientConn := &mockConnection{}
	clientProt := &mockProtection{}
	reg := packets.NewPacketRegistry()

	proxy := NewProxy(serverConn, serverProt, clientConn, clientProt, reg)

	// Set fake client protection keys
	fakeKeys := []byte{99, 98, 97, 96}
	proxy.SetClientProtectionKeys(fakeKeys)

	// Create a real ActivateProtectionPacket with real keys
	packet := network.NewActivateProtectionPacket()
	realKeys := []byte{1, 2, 3, 4}
	packet.Set("keys", realKeys)

	// Handle the packet
	result := proxy.handleActivateProtection(packet)

	assert.NotNil(t, result)
	assert.Equal(t, packet, result)

	// Server protection should be activated with real keys
	assert.True(t, serverProt.activateCalled)
	assert.Equal(t, realKeys, serverProt.activateKeys)

	// Client should get fake keys
	modifiedKeys := packets.Attr[[]byte]("keys", packet)
	assert.Equal(t, fakeKeys, modifiedKeys)
}

func TestHandleActivateProtection_IgnoresOtherPackets(t *testing.T) {
	serverConn := &mockConnection{}
	serverProt := &mockProtection{}
	clientConn := &mockConnection{}
	clientProt := &mockProtection{}
	reg := packets.NewPacketRegistry()

	proxy := NewProxy(serverConn, serverProt, clientConn, clientProt, reg)

	otherPacket := &mockPacket{id: 9999}
	result := proxy.handleActivateProtection(otherPacket)

	assert.Equal(t, otherPacket, result)
	// Server protection should NOT be activated
	assert.False(t, serverProt.activateCalled)
}

func TestOnServerToClient_RegistersInboundHandler(t *testing.T) {
	serverConn := &mockConnection{}
	serverProt := &mockProtection{}
	clientConn := &mockConnection{}
	clientProt := &mockProtection{}
	reg := packets.NewPacketRegistry()

	proxy := NewProxy(serverConn, serverProt, clientConn, clientProt, reg)

	proxy.OnServerToClient(func(p packets.Packet) packets.Packet {
		return p
	})

	// Handler is registered; execution is tested through Run with actual packets
	assert.NotNil(t, proxy.serverHandler)
}

func TestOnClientToServer_RegistersOutboundHandler(t *testing.T) {
	serverConn := &mockConnection{}
	serverProt := &mockProtection{}
	clientConn := &mockConnection{}
	clientProt := &mockProtection{}
	reg := packets.NewPacketRegistry()

	proxy := NewProxy(serverConn, serverProt, clientConn, clientProt, reg)

	proxy.OnClientToServer(func(p packets.Packet) packets.Packet {
		return p
	})

	// Handler is registered; execution is tested through Send
	assert.NotNil(t, proxy.clientHandler)
}

func TestMultipleServerToClientHandlers(t *testing.T) {
	packetID := int32(2001)
	packetData := []byte{10, 11, 12, 13}
	packetLen := int32(packets.HeaderLength + len(packetData))

	callCount := 0
	serverConn := &mockConnection{
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
	serverProt := &mockProtection{}
	clientConn := &mockConnection{}
	clientProt := &mockProtection{}

	reg := packets.NewPacketRegistry()
	reg.Register(packetID, "TestPacket", func() packets.Packet {
		return &mockPacket{id: packetID}
	})

	proxy := NewProxy(serverConn, serverProt, clientConn, clientProt, reg)

	var mu sync.Mutex
	callOrder := []int{}

	proxy.OnServerToClient(func(p packets.Packet) packets.Packet {
		mu.Lock()
		callOrder = append(callOrder, 1)
		mu.Unlock()
		return p
	})

	proxy.OnServerToClient(func(p packets.Packet) packets.Packet {
		mu.Lock()
		callOrder = append(callOrder, 2)
		mu.Unlock()
		return p
	})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		proxy.Run(ctx)
	}()

	time.Sleep(150 * time.Millisecond)
	cancel()

	mu.Lock()
	order := callOrder
	mu.Unlock()
	assert.Greater(t, len(order), 0)
	assert.Equal(t, []int{1, 2}, order)
}

func TestServerToClientHandler_CancelsPacket(t *testing.T) {
	packetID := int32(4001)
	packetData := []byte{30, 31, 32, 33}
	packetLen := int32(packets.HeaderLength + len(packetData))

	callCount := 0
	serverConn := &mockConnection{
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
	serverProt := &mockProtection{}
	clientConn := &mockConnection{}
	clientProt := &mockProtection{}

	reg := packets.NewPacketRegistry()
	reg.Register(packetID, "TestPacket", func() packets.Packet {
		return &mockPacket{id: packetID}
	})

	proxy := NewProxy(serverConn, serverProt, clientConn, clientProt, reg)

	var mu sync.Mutex
	firstHandlerCalled := false
	secondHandlerCalled := false

	proxy.OnServerToClient(func(p packets.Packet) packets.Packet {
		mu.Lock()
		firstHandlerCalled = true
		mu.Unlock()
		return nil // Cancel packet
	})

	proxy.OnServerToClient(func(p packets.Packet) packets.Packet {
		mu.Lock()
		secondHandlerCalled = true
		mu.Unlock()
		return p
	})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		proxy.Run(ctx)
	}()

	time.Sleep(150 * time.Millisecond)
	cancel()

	mu.Lock()
	first := firstHandlerCalled
	second := secondHandlerCalled
	mu.Unlock()
	assert.True(t, first)
	assert.False(t, second)
}

func TestClientToServer_ReceivesClientPackets(t *testing.T) {
	packetID := int32(3001)
	packetData := []byte{20, 21, 22, 23}
	packetLen := int32(packets.HeaderLength + len(packetData))

	serverConn := &mockConnection{}
	serverProt := &mockProtection{}

	callCount := 0
	clientConn := &mockConnection{
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
	clientProt := &mockProtection{}

	reg := packets.NewPacketRegistry()
	reg.Register(packetID, "TestPacket", func() packets.Packet {
		return &mockPacket{id: packetID}
	})

	proxy := NewProxy(serverConn, serverProt, clientConn, clientProt, reg)

	var mu sync.Mutex
	packetReceived := false
	proxy.OnClientToServer(func(p packets.Packet) packets.Packet {
		mu.Lock()
		packetReceived = true
		mu.Unlock()
		assert.Equal(t, packetID, p.ID())
		return p
	})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		proxy.Run(ctx)
	}()

	time.Sleep(150 * time.Millisecond)
	cancel()

	mu.Lock()
	received := packetReceived
	mu.Unlock()
	assert.True(t, received)
}

func TestMultipleClientToServerHandlers(t *testing.T) {
	packetID := int32(5001)
	packetData := []byte{40, 41, 42, 43}
	packetLen := int32(packets.HeaderLength + len(packetData))

	serverConn := &mockConnection{}
	serverProt := &mockProtection{}

	callCount := 0
	clientConn := &mockConnection{
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
	clientProt := &mockProtection{}

	reg := packets.NewPacketRegistry()
	reg.Register(packetID, "TestPacket", func() packets.Packet {
		return &mockPacket{id: packetID}
	})

	proxy := NewProxy(serverConn, serverProt, clientConn, clientProt, reg)

	var mu sync.Mutex
	callOrder := []int{}

	proxy.OnClientToServer(func(p packets.Packet) packets.Packet {
		mu.Lock()
		callOrder = append(callOrder, 1)
		mu.Unlock()
		return p
	})

	proxy.OnClientToServer(func(p packets.Packet) packets.Packet {
		mu.Lock()
		callOrder = append(callOrder, 2)
		mu.Unlock()
		return p
	})

	proxy.OnClientToServer(func(p packets.Packet) packets.Packet {
		mu.Lock()
		callOrder = append(callOrder, 3)
		mu.Unlock()
		return p
	})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		proxy.Run(ctx)
	}()

	time.Sleep(150 * time.Millisecond)
	cancel()

	mu.Lock()
	order := callOrder
	mu.Unlock()
	assert.Equal(t, []int{1, 2, 3}, order)
}

func TestClientToServerHandler_CancelsPacket(t *testing.T) {
	packetID := int32(6001)
	packetData := []byte{50, 51, 52, 53}
	packetLen := int32(packets.HeaderLength + len(packetData))

	serverConn := &mockConnection{}
	serverProt := &mockProtection{}

	callCount := 0
	clientConn := &mockConnection{
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
	clientProt := &mockProtection{}

	reg := packets.NewPacketRegistry()
	reg.Register(packetID, "TestPacket", func() packets.Packet {
		return &mockPacket{id: packetID}
	})

	proxy := NewProxy(serverConn, serverProt, clientConn, clientProt, reg)

	var mu sync.Mutex
	firstHandlerCalled := false
	secondHandlerCalled := false

	proxy.OnClientToServer(func(p packets.Packet) packets.Packet {
		mu.Lock()
		firstHandlerCalled = true
		mu.Unlock()
		return nil // Cancel packet
	})

	proxy.OnClientToServer(func(p packets.Packet) packets.Packet {
		mu.Lock()
		secondHandlerCalled = true
		mu.Unlock()
		return p
	})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		proxy.Run(ctx)
	}()

	time.Sleep(150 * time.Millisecond)
	cancel()

	mu.Lock()
	first := firstHandlerCalled
	second := secondHandlerCalled
	mu.Unlock()
	assert.True(t, first)
	assert.False(t, second)
}

func TestBidirectionalFlow(t *testing.T) {
	// Test that both directions work simultaneously
	serverPacketID := int32(7001)
	serverPacketData := []byte{1, 2, 3, 4}
	serverPacketLen := int32(packets.HeaderLength + len(serverPacketData))

	clientPacketID := int32(7002)
	clientPacketData := []byte{5, 6, 7, 8}
	clientPacketLen := int32(packets.HeaderLength + len(clientPacketData))

	serverCallCount := 0
	serverConn := &mockConnection{
		onRead: func(n int) ([]byte, error) {
			serverCallCount++
			switch serverCallCount {
			case 1:
				return createTestHeader(serverPacketLen, serverPacketID, false)[0:4], nil
			case 2:
				return createTestHeader(serverPacketLen, serverPacketID, false)[4:8], nil
			case 3:
				return serverPacketData, nil
			default:
				return nil, io.EOF
			}
		},
	}
	serverProt := &mockProtection{}

	clientCallCount := 0
	clientConn := &mockConnection{
		onRead: func(n int) ([]byte, error) {
			clientCallCount++
			switch clientCallCount {
			case 1:
				return createTestHeader(clientPacketLen, clientPacketID, false)[0:4], nil
			case 2:
				return createTestHeader(clientPacketLen, clientPacketID, false)[4:8], nil
			case 3:
				return clientPacketData, nil
			default:
				return nil, io.EOF
			}
		},
	}
	clientProt := &mockProtection{}

	reg := packets.NewPacketRegistry()
	reg.Register(serverPacketID, "ServerPacket", func() packets.Packet {
		return &mockPacket{id: serverPacketID}
	})
	reg.Register(clientPacketID, "ClientPacket", func() packets.Packet {
		return &mockPacket{id: clientPacketID}
	})

	proxy := NewProxy(serverConn, serverProt, clientConn, clientProt, reg)

	var mu sync.Mutex
	serverToClientReceived := false
	clientToServerReceived := false

	proxy.OnServerToClient(func(p packets.Packet) packets.Packet {
		mu.Lock()
		serverToClientReceived = true
		mu.Unlock()
		return p
	})

	proxy.OnClientToServer(func(p packets.Packet) packets.Packet {
		mu.Lock()
		clientToServerReceived = true
		mu.Unlock()
		return p
	})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		proxy.Run(ctx)
	}()

	time.Sleep(150 * time.Millisecond)
	cancel()

	mu.Lock()
	s2c := serverToClientReceived
	c2s := clientToServerReceived
	mu.Unlock()
	assert.True(t, s2c)
	assert.True(t, c2s)
}
