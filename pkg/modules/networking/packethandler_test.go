package networking

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/the-new-day/probogo/pkg/packets"
)

func TestHandlerSend_Success(t *testing.T) {
	mockConn := &mockConnection{}
	mockProt := &mockProtection{}
	reg := packets.NewPacketRegistry()

	handler := NewPacketHandler(mockConn, mockProt, reg)
	packet := &mockPacket{id: 1001}

	err := handler.Send(packet)

	assert.NoError(t, err)
	assert.Len(t, mockConn.writeData, 1)
}

func TestHandlerSend_WithOutboundHandler(t *testing.T) {
	mockConn := &mockConnection{}
	mockProt := &mockProtection{}
	reg := packets.NewPacketRegistry()

	handler := NewPacketHandler(mockConn, mockProt, reg)

	handlerCalled := false
	handler.OnOutBound(func(p packets.Packet) packets.Packet {
		handlerCalled = true
		return p
	})

	packet := &mockPacket{id: 1001}
	err := handler.Send(packet)

	assert.NoError(t, err)
	assert.True(t, handlerCalled)
	assert.Len(t, mockConn.writeData, 1)
}

func TestHandlerSend_OutboundHandlerModifiesPacket(t *testing.T) {
	mockConn := &mockConnection{}
	mockProt := &mockProtection{}
	reg := packets.NewPacketRegistry()

	handler := NewPacketHandler(mockConn, mockProt, reg)

	handler.OnOutBound(func(p packets.Packet) packets.Packet {
		// Return a different packet
		return &mockPacket{id: 2002}
	})

	originalPacket := &mockPacket{id: 1001}
	err := handler.Send(originalPacket)

	assert.NoError(t, err)
	// Packet should still be sent
	assert.Len(t, mockConn.writeData, 1)
}

func TestHandlerSend_OutboundHandlerCancelsPacket(t *testing.T) {
	mockConn := &mockConnection{}
	mockProt := &mockProtection{}
	reg := packets.NewPacketRegistry()

	handler := NewPacketHandler(mockConn, mockProt, reg)

	handler.OnOutBound(func(p packets.Packet) packets.Packet {
		return nil // Cancel packet
	})

	packet := &mockPacket{id: 1001}
	err := handler.Send(packet)

	// Should not return an error, but packet should not be sent
	assert.NoError(t, err)
	assert.Len(t, mockConn.writeData, 0)
}

func TestHandlerSend_MultipleOutboundHandlers(t *testing.T) {
	mockConn := &mockConnection{}
	mockProt := &mockProtection{}
	reg := packets.NewPacketRegistry()

	handler := NewPacketHandler(mockConn, mockProt, reg)

	callOrder := []int{}

	handler.OnOutBound(func(p packets.Packet) packets.Packet {
		callOrder = append(callOrder, 1)
		return p
	})

	handler.OnOutBound(func(p packets.Packet) packets.Packet {
		callOrder = append(callOrder, 2)
		return p
	})

	handler.OnOutBound(func(p packets.Packet) packets.Packet {
		callOrder = append(callOrder, 3)
		return p
	})

	packet := &mockPacket{id: 1001}
	err := handler.Send(packet)

	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, callOrder)
	assert.Len(t, mockConn.writeData, 1)
}

func TestHandlerSend_OutboundHandlerStopsChain(t *testing.T) {
	mockConn := &mockConnection{}
	mockProt := &mockProtection{}
	reg := packets.NewPacketRegistry()

	handler := NewPacketHandler(mockConn, mockProt, reg)

	callCount := 0

	handler.OnOutBound(func(p packets.Packet) packets.Packet {
		callCount++
		return p
	})

	handler.OnOutBound(func(p packets.Packet) packets.Packet {
		// Cancel in second handler
		return nil
	})

	handler.OnOutBound(func(p packets.Packet) packets.Packet {
		// This should not be called
		callCount++
		return p
	})

	packet := &mockPacket{id: 1001}
	err := handler.Send(packet)

	assert.NoError(t, err)
	assert.Equal(t, 1, callCount) // Only first handler called
	assert.Len(t, mockConn.writeData, 0)
}

func TestHandlerSend_PacketStreamError(t *testing.T) {
	expectedErr := errors.New("stream send failed")
	mockConn := &mockConnection{writeErr: expectedErr}
	mockProt := &mockProtection{}
	reg := packets.NewPacketRegistry()

	handler := NewPacketHandler(mockConn, mockProt, reg)
	packet := &mockPacket{id: 1001}

	err := handler.Send(packet)

	assert.ErrorIs(t, err, expectedErr)
}

func TestHandlerRun_ContextCancel(t *testing.T) {
	mockConn := &mockConnection{}
	mockProt := &mockProtection{}
	reg := packets.NewPacketRegistry()

	handler := NewPacketHandler(mockConn, mockProt, reg)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Should exit immediately
	done := make(chan bool)
	go func() {
		handler.Run(ctx)
		done <- true
	}()

	<-done // Should complete without hanging
}

func TestHandlerRun_InboundPacketSuccess(t *testing.T) {
	packetID := int32(1001)
	packetData := []byte{1, 2, 3, 4}
	packetLen := int32(packets.HeaderLength + len(packetData))

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
	handler := NewPacketHandler(mockConn, mockProt, reg)

	packetReceived := false
	handler.OnInBound(func(p packets.Packet) packets.Packet {
		packetReceived = true
		assert.Equal(t, packetID, p.ID())
		return p
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		handler.Run(ctx)
	}()

	// Allow time for the packet to be read
	time.Sleep(100 * time.Millisecond)
	cancel()

	assert.True(t, packetReceived)
}

func TestHandlerRun_InboundHandlerModifiesPacket(t *testing.T) {
	packetID := int32(1001)
	packetData := []byte{1, 2, 3, 4}
	packetLen := int32(packets.HeaderLength + len(packetData))

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
	handler := NewPacketHandler(mockConn, mockProt, reg)

	modifiedCalled := false
	handler.OnInBound(func(p packets.Packet) packets.Packet {
		// Handler can modify the packet
		modifiedCalled = true
		return p
	})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		handler.Run(ctx)
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()

	assert.True(t, modifiedCalled)
}

func TestHandlerRun_InboundHandlerCancelsPacket(t *testing.T) {
	packetID := int32(1001)
	packetData := []byte{1, 2, 3, 4}
	packetLen := int32(packets.HeaderLength + len(packetData))

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
	handler := NewPacketHandler(mockConn, mockProt, reg)

	firstHandlerCalled := false
	secondHandlerCalled := false

	handler.OnInBound(func(p packets.Packet) packets.Packet {
		firstHandlerCalled = true
		return nil // Cancel packet
	})

	handler.OnInBound(func(p packets.Packet) packets.Packet {
		secondHandlerCalled = true
		return p
	})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		handler.Run(ctx)
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()

	assert.True(t, firstHandlerCalled)
	assert.False(t, secondHandlerCalled)
}

func TestHandlerRun_MultipleInboundHandlers(t *testing.T) {
	packetID := int32(1001)
	packetData := []byte{1, 2, 3, 4}
	packetLen := int32(packets.HeaderLength + len(packetData))

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
	handler := NewPacketHandler(mockConn, mockProt, reg)

	callOrder := []int{}

	handler.OnInBound(func(p packets.Packet) packets.Packet {
		callOrder = append(callOrder, 1)
		return p
	})

	handler.OnInBound(func(p packets.Packet) packets.Packet {
		callOrder = append(callOrder, 2)
		return p
	})

	handler.OnInBound(func(p packets.Packet) packets.Packet {
		callOrder = append(callOrder, 3)
		return p
	})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		handler.Run(ctx)
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()

	assert.Equal(t, []int{1, 2, 3}, callOrder)
}

func TestHandlerRun_ReceiveError(t *testing.T) {
	expectedErr := errors.New("read error")
	mockConn := &mockConnection{
		onRead: func(n int) ([]byte, error) {
			return nil, expectedErr
		},
	}

	mockProt := &mockProtection{}
	reg := packets.NewPacketRegistry()

	handler := NewPacketHandler(mockConn, mockProt, reg)

	errorReceived := false
	handler.OnReceiveError(func(res PacketResult) {
		errorReceived = true
		assert.ErrorIs(t, res.Err, expectedErr)
		assert.Nil(t, res.Packet)
	})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		handler.Run(ctx)
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()

	assert.True(t, errorReceived)
}

func TestHandlerRun_MultipleErrorHandlers(t *testing.T) {
	expectedErr := errors.New("read error")
	callCount := 0
	mockConn := &mockConnection{
		onRead: func(n int) ([]byte, error) {
			callCount++
			if callCount <= 3 { // Return error on first read attempts
				return nil, expectedErr
			}
			return nil, io.EOF
		},
	}

	mockProt := &mockProtection{}
	reg := packets.NewPacketRegistry()

	handler := NewPacketHandler(mockConn, mockProt, reg)

	errorHandlerCount := 0

	handler.OnReceiveError(func(res PacketResult) {
		errorHandlerCount++
	})

	handler.OnReceiveError(func(res PacketResult) {
		assert.ErrorIs(t, res.Err, expectedErr)
	})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		handler.Run(ctx)
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()

	assert.Greater(t, errorHandlerCount, 0)
}

func TestHandlerRun_ContinuesAfterError(t *testing.T) {
	// First packet: causes error
	packetID1 := int32(1001)

	// Second packet: successful
	packetID2 := int32(1002)
	packetData2 := []byte{5, 6, 7, 8}
	packetLen2 := int32(packets.HeaderLength + len(packetData2))

	reg := packets.NewPacketRegistry()
	reg.Register(packetID1, "TestPacket1", func() packets.Packet {
		return &mockPacket{id: packetID1, unwrapErr: errors.New("unwrap failed")}
	})
	reg.Register(packetID2, "TestPacket2", func() packets.Packet {
		return &mockPacket{id: packetID2}
	})

	callCount := 0
	mockConn := &mockConnection{
		onRead: func(n int) ([]byte, error) {
			callCount++
			switch callCount {
			case 1:
				return createTestHeader(int32(packets.HeaderLength+1), packetID1, false)[0:4], nil
			case 2:
				return createTestHeader(int32(packets.HeaderLength+1), packetID1, false)[4:8], nil
			case 3:
				return []byte{0}, nil
			case 4:
				return createTestHeader(packetLen2, packetID2, false)[0:4], nil
			case 5:
				return createTestHeader(packetLen2, packetID2, false)[4:8], nil
			case 6:
				return packetData2, nil
			default:
				return nil, io.EOF
			}
		},
	}

	mockProt := &mockProtection{}
	handler := NewPacketHandler(mockConn, mockProt, reg)

	errorCount := 0
	packetCount := 0

	handler.OnReceiveError(func(res PacketResult) {
		errorCount++
	})

	handler.OnInBound(func(p packets.Packet) packets.Packet {
		packetCount++
		return p
	})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		handler.Run(ctx)
	}()

	time.Sleep(200 * time.Millisecond)
	cancel()

	assert.Greater(t, errorCount, 0)
	assert.Greater(t, packetCount, 0)
}

func TestHandlerActivateProtection(t *testing.T) {
	mockConn := &mockConnection{}
	mockProt := &mockProtection{}
	reg := packets.NewPacketRegistry()

	handler := NewPacketHandler(mockConn, mockProt, reg)
	keys := []byte{1, 2, 3, 4}

	handler.ActivateProtection(keys)

	assert.True(t, mockProt.activateCalled)
	assert.Equal(t, keys, mockProt.activateKeys)
}

func TestHandlerSend_WithNoHandlers(t *testing.T) {
	mockConn := &mockConnection{}
	mockProt := &mockProtection{}
	reg := packets.NewPacketRegistry()

	handler := NewPacketHandler(mockConn, mockProt, reg)
	packet := &mockPacket{id: 1001}

	err := handler.Send(packet)

	assert.NoError(t, err)
	assert.Len(t, mockConn.writeData, 1)
}

func TestHandlerRun_WithNoHandlers(t *testing.T) {
	packetID := int32(1001)
	packetData := []byte{1, 2, 3, 4}
	packetLen := int32(packets.HeaderLength + len(packetData))

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
	handler := NewPacketHandler(mockConn, mockProt, reg)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		handler.Run(ctx)
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()
}
