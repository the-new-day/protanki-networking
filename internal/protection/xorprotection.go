package protection

const vectorLength = 8

// Handles encryption and decryption of packet data using a custom XOR-based scheme.
type XorProtection struct {
	active           bool
	flipDirection    bool
	base             byte
	decryptionVector [vectorLength]byte
	encryptionVector [vectorLength]byte
	decryptionIndex  int
	encryptionIndex  int
}

// Creates a new instance of XorProtection.
//
// flipDirection - whether to flip the encryption/decryption roles
func NewXorProtection(flipDirection bool) *XorProtection {
	return &XorProtection{
		flipDirection: flipDirection,
	}
}

// Activates protection using a list of keys.
func (p *XorProtection) Activate(keys []byte) {
	for _, key := range keys {
		p.base ^= key
	}

	for i := range vectorLength {
		baseXor := p.base ^ byte(i<<3)

		if !p.flipDirection {
			p.decryptionVector[i] = baseXor
			p.encryptionVector[i] = baseXor ^ 0x57
		} else {
			p.decryptionVector[i] = baseXor ^ 0x57
			p.encryptionVector[i] = baseXor
		}
	}
	p.active = true
}

func (p *XorProtection) Encrypt(rawData []byte) []byte {
	result := make([]byte, len(rawData))
	copy(result, rawData)

	if !p.active || len(rawData) == 0 {
		return result
	}

	for i, rawByte := range rawData {
		result[i] = rawByte ^ p.encryptionVector[p.encryptionIndex]
		p.encryptionVector[p.encryptionIndex] = rawByte
		p.encryptionIndex ^= int(rawByte & 7)
	}

	return result
}

func (p *XorProtection) Decrypt(encryptedData []byte) []byte {
	result := make([]byte, len(encryptedData))
	copy(result, encryptedData)

	if !p.active || len(encryptedData) == 0 {
		return result
	}

	for i, encryptedByte := range encryptedData {
		decVal := encryptedByte ^ p.decryptionVector[p.decryptionIndex]
		p.decryptionVector[p.decryptionIndex] = decVal
		result[i] = decVal
		p.decryptionIndex ^= int(decVal & 7)
	}

	return result
}
