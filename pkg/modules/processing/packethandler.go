package processing

import (
	"fmt"

	"github.com/the-new-day/probogo/pkg/modules/networking"
	"github.com/the-new-day/probogo/pkg/modules/protection"
	"github.com/the-new-day/probogo/pkg/packets"
)

// PacketHandler denotes link in a chain of packet handlers.
// It allows build linear chains as well as hierarchical structures.
type PacketHandler interface {
	// Handle handles a packet. The chain link can modify the packet or react in any other way.
	// One chain link can stop iterating through the chain or call the next chain link.
	Handle(packet packets.Packet) (handled packets.Packet, err error)

	// AddNext adds the next chain link to the end of the chain.
	AddNext(handler PacketHandler)

	// Add adds child link to the chain, making it a tree.
	AddChild(handler PacketHandler)
}

// BasePacketHandler provides chaining functionality for all packet handlers.
// PacketHandlers chain is effectively a linked list, and BasePacketHandler provides
// functionality that allows "iterating" through that list.
type BasePacketHandler struct {
	next       PacketHandler
	child      PacketHandler
	socket     *networking.GameSocket
	protection protection.Protection
}

// NewBasePacketHandler creates new BasePacketHandler with next chain link set to nil.
func NewBasePacketHandler(socket *networking.GameSocket, protection protection.Protection) *BasePacketHandler {
	return &BasePacketHandler{nil, nil, socket, protection}
}

func (h *BasePacketHandler) AddNext(handler PacketHandler) {
	if h.next == nil {
		h.next = handler
	} else {
		h.next.AddNext(handler)
	}
}

func (h *BasePacketHandler) AddChild(handler PacketHandler) {
	if h.child == nil {
		h.child = handler
	} else {
		h.child.AddNext(handler)
	}
}

// CallNext passes the packet to the next handler in the chain.
// If there's no next handler, returns nil and an error.
func (h *BasePacketHandler) CallNext(packet packets.Packet) (handled packets.Packet, err error) {
	if h.next == nil {
		return nil, fmt.Errorf("BasePacketHandler.CallNext: missing terminal link")
	}
	return h.next.Handle(packet)
}

// CallNext passes the packet to the child handler of the link.
// If there's no child handler, returns nil and an error.
func (h *BasePacketHandler) CallChild(packet packets.Packet) (handled packets.Packet, err error) {
	if h.child == nil {
		return nil, fmt.Errorf("BasePacketHandler.CallChild: missing terminal link")
	}
	return h.child.Handle(packet)
}

// SendPacket sends packet through the socket.
func (h *BasePacketHandler) SendPacket(packet packets.Packet) error {
	wrapped, err := packet.Wrap(h.protection)
	if err != nil {
		return err
	}
	return h.socket.Send(wrapped.Bytes())
}

// SendPackets sends packets through the socket.
func (h *BasePacketHandler) SendPackets(packets []packets.Packet) error {
	packetsData := make([]byte, 0)

	for _, packet := range packets {
		wrapped, err := packet.Wrap(h.protection)
		if err != nil {
			return err
		}
		packetsData = append(packetsData, wrapped.Bytes()...)
	}

	return h.socket.Send(packetsData)
}

// CloseSocket closes current socket.
// It should only be called in the terminal link.
func (h *BasePacketHandler) CloseSocket() error {
	return h.socket.Close()
}
