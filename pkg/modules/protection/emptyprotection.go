package protection

type EmptyProtection struct{}

func (p *EmptyProtection) Activate([]byte) {}

func (p *EmptyProtection) Encrypt(data []byte) []byte {
	return data
}

func (p *EmptyProtection) Decrypt(data []byte) []byte {
	return data
}
