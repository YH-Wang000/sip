package locate

// StaticLocator Implement service discovery through manual advance configuration
type StaticLocator struct {
	registerServerAddress string
	proxyServerAddress    string
}

func (s *StaticLocator) LocateRegistrar(domain string) (string, error) {
	return s.registerServerAddress, nil
}

func (s *StaticLocator) LocateProxy(domain string) (string, error) {
	return s.proxyServerAddress, nil
}
