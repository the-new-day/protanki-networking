package protection

// Basic interface for protection.
type Protection interface {
	Encrypt(rawData []byte) []byte
	Decrypt(encryptedData []byte) []byte
}
