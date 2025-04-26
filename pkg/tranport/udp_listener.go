package tranport

import (
	"net"

	"sip/pkg/sipmsg"
)

type udpSocketListener struct {
	conn  net.Conn
	index TransportIndex
}

func NewUdpSocketListener(conn net.Conn) SocketListener {
	return &udpSocketListener{
		conn: conn,
	}
}

func (l *udpSocketListener) ReadMessage() (*sipmsg.GenericMessage, error) {
	//TODO implement me
	panic("implement me")
}

func (l *udpSocketListener) Start() {
	//TODO implement me
	panic("implement me")
}

func (l *udpSocketListener) Close() {
	//TODO implement me
	panic("implement me")
}
