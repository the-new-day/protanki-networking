package processing

import (
	"fmt"

	"github.com/the-new-day/probogo/pkg/packets"
)

// PacketHandler denotes link in a chain of packet handlers.
type PacketHandler interface {
	// Handle handles a packet. The chain link can modify the packet or react in any other way.
	// One chain link can stop iterating through the chain or call the next chain link.
	Handle(packet packets.Packet) (handled packets.Packet, err error)

	// AddNext adds the next chain link to the end of the chain.
	AddNext(handler PacketHandler)
}

// BasePacketHandler provides chaining functionality for all packet handlers.
// PacketHandlers chain is effectively a linked list, and BasePacketHandler provides
// functionality that allows "iterating" through that list.
type BasePacketHandler struct {
	next PacketHandler
}

// NewBasePacketHandler creates new BasePacketHandler with next chain link set to nil.
func NewBasePacketHandler() *BasePacketHandler {
	return &BasePacketHandler{nil}
}

// AddNext adds new link to the end of the chain.
func (h *BasePacketHandler) AddNext(handler PacketHandler) {
	if h.next == nil {
		h.next = handler
	} else {
		h.next.AddNext(handler)
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
