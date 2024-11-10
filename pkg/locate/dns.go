package locate

// DnsSrvRecordLocator Implement service discovery by querying dns srv records
type DnsSrvRecordLocator struct {
}

func (d *DnsSrvRecordLocator) LocateRegistrar(domain string) (string, error) {
	return "", nil
}

func (d *DnsSrvRecordLocator) LocateProxy(domain string) (string, error) {
	return "", nil
}
