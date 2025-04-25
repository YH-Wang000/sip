package tranport

import (
	"sip/pkg/sipmsg"
)

type MessageHandler func(msg *sipmsg.GenericMessage)

type SipClientTransport interface {
	SendRequest(request *sipmsg.GenericMessage) error
	SetResponseHandler(handler func(response *sipmsg.GenericMessage))
}

type clientTransportImpl struct {
	sipTransport SipTransport
}

func NewClientTransport() SipClientTransport {
	return nil
}

func (c *clientTransportImpl) Start() {
}

func (c *clientTransportImpl) SendRequest(request *sipmsg.GenericMessage) error {
	return c.sipTransport.SendMessage(request)
}

func (c *clientTransportImpl) SetResponseHandler(handler func(response *sipmsg.GenericMessage)) {
	if handler == nil {
		return
	}
	c.sipTransport.SetMessageHandler(func(msg *sipmsg.GenericMessage) {
		if msg.StartLine.IsStatusLine() {
			handler(msg)
		}
	})
}
