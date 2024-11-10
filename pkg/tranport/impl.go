package tranport

import (
	"sip/pkg/locate"
	"sip/pkg/sipmsg"
)

type transport struct {
	locator locate.ServerLocator
}

func NewTransport(locator locate.ServerLocator) Transport {
	return &transport{
		locator: locator,
	}
}

func (t *transport) SendRequest(request *sipmsg.GenericMessage) (ResponseFuture, error) {
	return nil, nil
}

func (t *transport) SendResponse(response *sipmsg.GenericMessage) error {
	return nil
}

func (t *transport) RegisterMessageHandler(f func(*sipmsg.GenericMessage)) {
}
