package networking

// Address represents a network address
type Address struct {
	Host     string
	Port     uint16
	Username string // for proxy auth
	Password string // for proxy auth
}

func NewAddress(host string, port uint16) *Address {
	return &Address{
		Host: host,
		Port: port,
	}
}

func NewAddressProxy(host string, port uint16, username string, password string) *Address {
	return &Address{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}
