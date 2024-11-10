package locate

// ServerLocator Used to determine server address and perform service discovery
// It can be configured manually, statically, through DNS SRV, or through DHCP
type ServerLocator interface {
	LocateRegistrar(domain string) (string, error)
	LocateProxy(domain string) (string, error)
}
