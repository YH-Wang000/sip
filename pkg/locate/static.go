package locate

// StaticLocator Implement service discovery through manual advance configuration
type StaticLocator struct {
	RegisterServerAddress string
	ProxyServerAddress    string
}

func (s *StaticLocator) LocateRegistrar(domain string) (string, error) {
	return s.RegisterServerAddress, nil
}

func (s *StaticLocator) LocateProxy(domain string) (string, error) {
	return s.ProxyServerAddress, nil
}
