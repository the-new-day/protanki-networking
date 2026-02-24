package protection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientServerInteraction(t *testing.T) {
	keys := []byte{0x12, 0x34, 0x56, 0x78}
	originalData := []byte("Hello World!")

	client := NewProtection(false)
	client.Activate(keys)

	server := NewProtection(true)
	server.Activate(keys)

	encrypted := client.Encrypt(originalData)

	assert.NotEqual(t, originalData, encrypted, "Encrypted data should not match original data")

	decrypted := server.Decrypt(encrypted)

	assert.Equal(t, originalData, decrypted, "Decryption failed")
}

func TestInactiveState(t *testing.T) {
	p := NewProtection(false)
	data := []byte{1, 2, 3, 4, 5}

	encrypted := p.Encrypt(data)
	assert.Equal(t, data, encrypted, "Inactive protection should not alter data on Encrypt")

	decrypted := p.Decrypt(data)
	assert.Equal(t, data, decrypted, "Inactive protection should not alter data on Decrypt")
}

func TestStatefulness(t *testing.T) {
	keys := []byte{0xAA}
	data := []byte{0x00, 0x00, 0x00, 0x00}

	p := NewProtection(false)
	p.Activate(keys)

	firstPass := p.Encrypt(data)
	secondPass := p.Encrypt(data)

	assert.NotEqual(t, firstPass, secondPass, "Stateful cipher should produce different outputs for repeated identical inputs")
}

func TestNoMutation(t *testing.T) {
	keys := []byte{0xFF}
	originalData := []byte{10, 20, 30}

	snapshot := make([]byte, len(originalData))
	copy(snapshot, originalData)

	p := NewProtection(false)
	p.Activate(keys)

	_ = p.Encrypt(originalData)

	assert.Equal(t, snapshot, originalData, "Encrypt mutated the original slice")
}
