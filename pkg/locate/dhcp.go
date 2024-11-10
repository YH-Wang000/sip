package locate

// DhcpLocator Automatically discover servers through dhcp protocol
type DhcpLocator struct {
}

func (d *DhcpLocator) LocateRegistrar(domain string) (string, error) {
	return "", nil
}

func (d *DhcpLocator) LocateProxy(domain string) (string, error) {
	return "", nil
}
