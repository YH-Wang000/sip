package tranport

import "sip/pkg/sipmsg"

type ResponseFuture interface {
	GetResponse() (*sipmsg.GenericMessage, error)
}

type ClientTransport interface {
	SendRequest(request *sipmsg.GenericMessage) (ResponseFuture, error)
}

type ServerTransport interface {
	SendResponse(response *sipmsg.GenericMessage) error
	RegisterMessageHandler(func(*sipmsg.GenericMessage))
}

type Transport interface {
	ClientTransport
	ServerTransport
}
