package packets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockPacket struct {
	BasePacket
}

func NewMockPacket() Packet {
	return &MockPacket{}
}

func TestNewPacketRegistry(t *testing.T) {
	registry := NewPacketRegistry()
	assert.NotNil(t, registry)
	assert.NotNil(t, registry.packets)
	assert.Empty(t, registry.packets)
}

func TestPacketRegistry_Register(t *testing.T) {
	registry := NewPacketRegistry()

	registry.Register(1001, "TestPacket", NewMockPacket)

	stored, exists := registry.packets[1001]
	assert.True(t, exists)
	assert.Equal(t, "TestPacket", stored.name)
	assert.NotNil(t, stored.factory)
}

func TestPacketRegistry_Register_DuplicatePanics(t *testing.T) {
	registry := NewPacketRegistry()

	registry.Register(1001, "FirstPacket", NewMockPacket)

	assert.Panics(t, func() {
		registry.Register(1001, "SecondPacket", NewMockPacket)
	})
}

func TestPacketRegistry_Get_Existing(t *testing.T) {
	registry := NewPacketRegistry()

	registry.Register(1001, "TestPacket", NewMockPacket)

	packet := registry.Get(1001)
	assert.NotNil(t, packet)
	assert.IsType(t, &MockPacket{}, packet)
}

func TestPacketRegistry_Get_NonExisting(t *testing.T) {
	registry := NewPacketRegistry()

	packet := registry.Get(9999)
	assert.Nil(t, packet)
}

func TestPacketRegistry_GetName_Existing(t *testing.T) {
	registry := NewPacketRegistry()

	registry.Register(1001, "TestPacket", NewMockPacket)

	name := registry.GetName(1001)
	assert.Equal(t, "TestPacket", name)
}

func TestPacketRegistry_GetName_NonExisting(t *testing.T) {
	registry := NewPacketRegistry()

	name := registry.GetName(9999)
	assert.Empty(t, name)
}

func TestPacketRegistry_MultiplePackets(t *testing.T) {
	registry := NewPacketRegistry()

	registry.Register(1001, "FirstPacket", NewMockPacket)
	registry.Register(1002, "SecondPacket", NewMockPacket)
	registry.Register(1003, "ThirdPacket", NewMockPacket)

	assert.Len(t, registry.packets, 3)

	p1 := registry.Get(1001)
	p2 := registry.Get(1002)
	p3 := registry.Get(1003)

	assert.NotNil(t, p1)
	assert.NotNil(t, p2)
	assert.NotNil(t, p3)

	assert.Equal(t, "FirstPacket", registry.GetName(1001))
	assert.Equal(t, "SecondPacket", registry.GetName(1002))
	assert.Equal(t, "ThirdPacket", registry.GetName(1003))
}

func TestPacketRegistry_FactoryReturnsNewInstance(t *testing.T) {
	registry := NewPacketRegistry()

	registry.Register(1001, "TestPacket", NewMockPacket)

	p1 := registry.Get(1001)
	p2 := registry.Get(1001)

	// Should be different instances
	assert.NotSame(t, p1, p2)
}
