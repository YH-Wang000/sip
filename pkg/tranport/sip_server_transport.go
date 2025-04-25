package tranport

import "sip/pkg/sipmsg"

type SipServerTransport interface {
	SendResponse(response *sipmsg.GenericMessage) error
	SetRequestHandler(handler func(request *sipmsg.GenericMessage))
}

func NewServerTransport() SipServerTransport {
	return &serverTransportImpl{}
}

type serverTransportImpl struct {
	sipTransport SipTransport
}

func (s *serverTransportImpl) SendResponse(response *sipmsg.GenericMessage) error {
	return s.sipTransport.SendMessage(response)
}

func (s *serverTransportImpl) SetRequestHandler(handler func(request *sipmsg.GenericMessage)) {
	if handler == nil {
		return
	}
	s.sipTransport.SetMessageHandler(func(msg *sipmsg.GenericMessage) {
		if msg.StartLine.IsRequestLine() {
			handler(msg)
		}
	})
}
