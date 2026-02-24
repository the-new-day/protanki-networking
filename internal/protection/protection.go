package protection

const VectorLength = 8

type Protection struct {
	active           bool
	flipDirection    bool
	base             byte
	decryptionVector [VectorLength]byte
	encryptionVector [VectorLength]byte
	decryptionIndex  int
	encryptionIndex  int
}

func NewProtection(flipDirection bool) *Protection {
	return &Protection{
		flipDirection: flipDirection,
	}
}

func (p *Protection) Activate(keys []byte) {
	for _, key := range keys {
		p.base ^= key
	}

	for i := range VectorLength {
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

func (p *Protection) Encrypt(rawData []byte) []byte {
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

func (p *Protection) Decrypt(encryptedData []byte) []byte {
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
