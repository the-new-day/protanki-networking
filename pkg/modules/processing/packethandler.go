package processing

import (
	"context"
	"fmt"

	"github.com/the-new-day/probogo/pkg/packets"
)

// PacketHandler denotes link in a chain of packet handlers.
// It allows build linear chains as well as hierarchical structures.
type PacketHandler interface {
	// Handle handles a packet. The chain link can modify the packet or react in any other way.
	// One chain link can stop iterating through the chain or call the next chain link.
	Handle(packet packets.Packet, ctx *context.Context) (handled packets.Packet, err error)

	// AddNext adds the next chain link to the end of the chain.
	AddNext(handler PacketHandler)

	// Add adds child link to the chain, making it a tree.
	AddChild(handler PacketHandler)
}

// BasePacketHandler provides chaining functionality for all packet handlers.
// PacketHandlers chain is effectively a linked list, and BasePacketHandler provides
// functionality that allows "iterating" through that list.
type BasePacketHandler struct {
	next  PacketHandler
	child PacketHandler
}

// NewBasePacketHandler creates new BasePacketHandler with next chain link set to nil.
func NewBasePacketHandler() *BasePacketHandler {
	return &BasePacketHandler{nil, nil}
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
func (h *BasePacketHandler) CallNext(packet packets.Packet, ctx *context.Context) (handled packets.Packet, err error) {
	if h.next == nil {
		return nil, fmt.Errorf("BasePacketHandler.CallNext: missing terminal link")
	}
	return h.next.Handle(packet, ctx)
}

// CallNext passes the packet to the child handler of the link.
// If there's no child handler, returns nil and an error.
func (h *BasePacketHandler) CallChild(packet packets.Packet, ctx *context.Context) (handled packets.Packet, err error) {
	if h.child == nil {
		return nil, fmt.Errorf("BasePacketHandler.CallChild: missing terminal link")
	}
	return h.child.Handle(packet, ctx)
}
