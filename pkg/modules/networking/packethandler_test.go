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

	handlerCalled := make(chan struct{})
	handler.OnOutBound(func(p packets.Packet) packets.Packet {
		close(handlerCalled)
		return p
	})

	packet := &mockPacket{id: 1001}
	err := handler.Send(packet)

	assert.NoError(t, err)
	select {
	case <-handlerCalled:
		// success
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for handler")
	}
	assert.Len(t, mockConn.writeData, 1)
}

func TestHandlerSend_OutboundHandlerModifiesPacket(t *testing.T) {
	mockConn := &mockConnection{}
	mockProt := &mockProtection{}
	reg := packets.NewPacketRegistry()

	handler := NewPacketHandler(mockConn, mockProt, reg)

	sentPacket := make(chan packets.Packet, 1)
	handler.OnOutBound(func(p packets.Packet) packets.Packet {
		modified := &mockPacket{id: 2002}
		sentPacket <- modified
		return modified
	})

	originalPacket := &mockPacket{id: 1001}
	err := handler.Send(originalPacket)

	assert.NoError(t, err)
	select {
	case pkt := <-sentPacket:
		assert.Equal(t, int32(2002), pkt.ID())
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for packet")
	}
	assert.Len(t, mockConn.writeData, 1)
}

func TestHandlerSend_OutboundHandlerCancelsPacket(t *testing.T) {
	mockConn := &mockConnection{}
	mockProt := &mockProtection{}
	reg := packets.NewPacketRegistry()

	handler := NewPacketHandler(mockConn, mockProt, reg)

	handlerCalled := make(chan struct{})
	handler.OnOutBound(func(p packets.Packet) packets.Packet {
		close(handlerCalled)
		return nil
	})

	packet := &mockPacket{id: 1001}
	err := handler.Send(packet)

	assert.NoError(t, err)
	select {
	case <-handlerCalled:
		// success
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for handler")
	}
	assert.Len(t, mockConn.writeData, 0)
}

func TestHandlerSend_MultipleOutboundHandlers(t *testing.T) {
	mockConn := &mockConnection{}
	mockProt := &mockProtection{}
	reg := packets.NewPacketRegistry()

	handler := NewPacketHandler(mockConn, mockProt, reg)

	callOrder := make(chan int, 3)
	done := make(chan struct{})

	handler.OnOutBound(func(p packets.Packet) packets.Packet {
		callOrder <- 1
		return p
	})

	handler.OnOutBound(func(p packets.Packet) packets.Packet {
		callOrder <- 2
		return p
	})

	handler.OnOutBound(func(p packets.Packet) packets.Packet {
		callOrder <- 3
		close(done)
		return p
	})

	packet := &mockPacket{id: 1001}
	err := handler.Send(packet)

	assert.NoError(t, err)

	select {
	case <-done:
		close(callOrder)
		collected := make([]int, 0, 3)
		for v := range callOrder {
			collected = append(collected, v)
		}
		assert.Equal(t, []int{1, 2, 3}, collected)
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for handlers")
	}
	assert.Len(t, mockConn.writeData, 1)
}

func TestHandlerSend_OutboundHandlerStopsChain(t *testing.T) {
	mockConn := &mockConnection{}
	mockProt := &mockProtection{}
	reg := packets.NewPacketRegistry()

	handler := NewPacketHandler(mockConn, mockProt, reg)

	firstCalled := make(chan struct{})
	secondCalled := make(chan struct{})

	handler.OnOutBound(func(p packets.Packet) packets.Packet {
		close(firstCalled)
		return p
	})

	handler.OnOutBound(func(p packets.Packet) packets.Packet {
		close(secondCalled)
		return nil
	})

	handler.OnOutBound(func(p packets.Packet) packets.Packet {
		t.Error("third handler should not be called")
		return p
	})

	packet := &mockPacket{id: 1001}
	err := handler.Send(packet)

	assert.NoError(t, err)

	select {
	case <-firstCalled:
		// ok
	case <-time.After(time.Second):
		t.Fatal("first handler not called")
	}

	select {
	case <-secondCalled:
		// ok
	case <-time.After(time.Second):
		t.Fatal("second handler not called")
	}

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

	done := make(chan struct{})
	go func() {
		handler.Run(ctx)
		close(done)
	}()

	select {
	case <-done:
		// ok
	case <-time.After(time.Second):
		t.Fatal("Run didn't exit after context cancel")
	}
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

	packetReceived := make(chan struct{})
	handler.OnInBound(func(p packets.Packet) packets.Packet {
		assert.Equal(t, packetID, p.ID())
		close(packetReceived)
		return p
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go handler.Run(ctx)

	select {
	case <-packetReceived:
		// success
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for packet")
	}
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

	modifiedCalled := make(chan struct{})
	handler.OnInBound(func(p packets.Packet) packets.Packet {
		close(modifiedCalled)
		return p
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go handler.Run(ctx)

	select {
	case <-modifiedCalled:
		// ok
	case <-time.After(time.Second):
		t.Fatal("handler not called")
	}
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

	firstCalled := make(chan struct{})
	secondCalled := make(chan struct{})

	handler.OnInBound(func(p packets.Packet) packets.Packet {
		close(firstCalled)
		return nil
	})

	handler.OnInBound(func(p packets.Packet) packets.Packet {
		close(secondCalled)
		return p
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go handler.Run(ctx)

	select {
	case <-firstCalled:
		// ok
	case <-time.After(time.Second):
		t.Fatal("first handler not called")
	}

	select {
	case <-secondCalled:
		t.Fatal("second handler should not be called")
	case <-time.After(100 * time.Millisecond):
		// ok
	}
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

	callOrder := make(chan int, 3)
	done := make(chan struct{})

	handler.OnInBound(func(p packets.Packet) packets.Packet {
		callOrder <- 1
		return p
	})

	handler.OnInBound(func(p packets.Packet) packets.Packet {
		callOrder <- 2
		return p
	})

	handler.OnInBound(func(p packets.Packet) packets.Packet {
		callOrder <- 3
		close(done)
		return p
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go handler.Run(ctx)

	select {
	case <-done:
		close(callOrder)
		collected := make([]int, 0, 3)
		for v := range callOrder {
			collected = append(collected, v)
		}
		assert.Equal(t, []int{1, 2, 3}, collected)
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for handlers")
	}
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

	errorReceived := make(chan PacketResult, 1)
	handler.OnError(func(res PacketResult) {
		errorReceived <- res
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go handler.Run(ctx)

	select {
	case res := <-errorReceived:
		assert.ErrorIs(t, res.Err, expectedErr)
		assert.Nil(t, res.Packet)
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for error handler")
	}
}

func TestHandlerRun_MultipleErrorHandlers(t *testing.T) {
	expectedErr := errors.New("read error")
	mockConn := &mockConnection{
		onRead: func(n int) ([]byte, error) {
			return nil, expectedErr
		},
	}

	mockProt := &mockProtection{}
	reg := packets.NewPacketRegistry()

	handler := NewPacketHandler(mockConn, mockProt, reg)

	firstDone := make(chan struct{})
	secondDone := make(chan struct{})

	handler.OnError(func(res PacketResult) {
		close(firstDone)
	})

	handler.OnError(func(res PacketResult) {
		close(secondDone)
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go handler.Run(ctx)

	for range 2 {
		select {
		case <-firstDone:
			firstDone = nil
		case <-secondDone:
			secondDone = nil
		case <-time.After(time.Second):
			t.Fatal("timeout waiting for error handlers")
		}
	}
}

func TestHandlerRun_ContinuesAfterError(t *testing.T) {
	packetID := int32(1002)
	packetData2 := []byte{5, 6, 7, 8}
	packetLen2 := int32(packets.HeaderLength + len(packetData2))

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
				// first read fails
				return nil, errors.New("temporary error")
			case 2:
				return createTestHeader(packetLen2, packetID, false)[0:4], nil
			case 3:
				return createTestHeader(packetLen2, packetID, false)[4:8], nil
			case 4:
				return packetData2, nil
			default:
				return nil, io.EOF
			}
		},
	}

	mockProt := &mockProtection{}
	handler := NewPacketHandler(mockConn, mockProt, reg)

	errorCount := 0
	errorDone := make(chan struct{})
	packetDone := make(chan struct{})

	handler.OnError(func(res PacketResult) {
		errorCount++
		if errorCount == 1 {
			close(errorDone)
		}
	})

	handler.OnInBound(func(p packets.Packet) packets.Packet {
		close(packetDone)
		return p
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go handler.Run(ctx)

	select {
	case <-errorDone:
		// ok
	case <-time.After(time.Second):
		t.Fatal("error handler not called")
	}

	select {
	case <-packetDone:
		t.Fatal("packet handler is called after error")
	default:
	}

	assert.Equal(t, 1, errorCount)
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
	done := make(chan struct{})
	go func() {
		handler.Run(ctx)
		close(done)
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()

	select {
	case <-done:
		// ok
	case <-time.After(time.Second):
		t.Fatal("Run didn't exit after cancel")
	}
}
