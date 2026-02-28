package networking

// Address represents a network address
type Address struct {
	Host string
	Port uint16
}

func NewAddress(host string, port uint16) *Address {
	return &Address{
		Host: host,
		Port: port,
	}
}
