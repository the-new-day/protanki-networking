package packets

import "fmt"

// Global packet registry.
var registry = NewPacketRegistry()

// Register registers new packet factory to the global registry.
// It panics if a factory with specified ID already exists.
func Register(id int32, factory PacketFactory) {
	registry.Register(id, factory)
}

// GetPacketByID returns packet produced by provided
// packet factory for the corresponding ID from the global registry.
// It return nil if a factory with specified ID wasn't previously registered.
func GetPacketByID(id int32) Packet {
	return registry.Get(id)
}

// PacketFactory is a function that returns initialized packet object.
type PacketFactory func() Packet

// PacketRegistry stores all packets with their factories (or constructors).
//
// Networking modules use this registry to find packets by their ID,
// so any new packets should be registered in the global registry using packets.Registry()
//
// It does not provide thread-safety as registering packets should occur only during initialization.
type PacketRegistry struct {
	factories map[int32]PacketFactory
}

func NewPacketRegistry() *PacketRegistry {
	return &PacketRegistry{factories: make(map[int32]PacketFactory)}
}

// Register registers new packet factory.
// It panics if a factory with specified id already exists.
func (pr *PacketRegistry) Register(id int32, factory PacketFactory) {
	if _, exists := pr.factories[id]; exists {
		panic(fmt.Sprintf("PacketRegistry: can't register packet %q: already exists", id))
	}
	pr.factories[id] = factory
}

// Get returns packet produced by provided packet factory for the corresponding id.
// It return nil if a factory with specified id wasn't previously registered.
func (pr *PacketRegistry) Get(id int32) Packet {
	factory, ok := pr.factories[id]
	if !ok {
		return nil
	}
	return factory()
}
