package packets

import "fmt"

// Global packet registry.
var registry = NewPacketRegistry()

// Register registers new packet factory to the global registry.
// It panics if a factory with specified ID already exists.
func Register(id int32, name string, factory PacketFactory) {
	registry.Register(id, name, factory)
}

// Get returns packet produced by provided
// packet factory for the corresponding ID from the global registry.
// It return nil if a factory with specified ID wasn't previously registered.
func Get(id int32) Packet {
	return registry.Get(id)
}

// GetName returns packet name for the corresponding ID.
// It returns an empty string if a packet with specified ID wasn't previously registered.
func GetName(id int32) string {
	return registry.GetName(id)
}

// PacketFactory is a function that returns initialized packet object.
type PacketFactory func() Packet

type storedPacket struct {
	factory PacketFactory
	name    string
}

// PacketRegistry stores all packets with their factories (or constructors).
//
// Networking modules use this registry to find packets by their ID,
// so any new packets should be registered in the global registry using packets.Registry()
//
// It does not provide thread-safety as registering packets should occur only during initialization.
type PacketRegistry struct {
	packets map[int32]storedPacket
}

func NewPacketRegistry() *PacketRegistry {
	return &PacketRegistry{packets: make(map[int32]storedPacket)}
}

// Register registers new packet factory.
// It panics if a factory with specified id already exists.
func (pr *PacketRegistry) Register(id int32, name string, factory PacketFactory) {
	if _, exists := pr.packets[id]; exists {
		panic(fmt.Sprintf("PacketRegistry: can't register packet %q (%q): already exists", id, name))
	}
	pr.packets[id] = storedPacket{
		factory: factory,
		name:    name,
	}

}

// Get returns packet produced by provided packet factory for the corresponding id.
// It return nil if a factory with specified id wasn't previously registered.
func (pr *PacketRegistry) Get(id int32) Packet {
	stored, ok := pr.packets[id]
	if !ok {
		return nil
	}
	return stored.factory()
}

// GetName returns packet name for the corresponding ID.
// It returns an empty string if a packet with specified ID wasn't previously registered.
func (pr *PacketRegistry) GetName(id int32) string {
	stored, ok := pr.packets[id]
	if !ok {
		return ""
	}
	return stored.name
}
