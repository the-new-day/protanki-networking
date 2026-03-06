package protection

// Basic interface for protection.
type Protection interface {
	// Activate activates protection, i.e. sets internal state, using a sequence of keys (bytes).
	// Resets internal state before activating.
	Activate(keys []byte)

	// Encrypt encrypts sequence of bytes. It may change internal state.
	Encrypt(rawData []byte) []byte

	// Decrypt decrypts sequence of bytes. It may change internal state.
	Decrypt(encryptedData []byte) []byte
}
