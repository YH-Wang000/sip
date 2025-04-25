package tranport

import (
	"net"

	"sip/pkg/enc"
	"sip/pkg/sipmsg"
)

type UdpMessageSender struct {
	sipTransport *sipTransportImpl
	index        TransportIndex
	encoder      enc.Encoder
}

func NewUdpMessageSender(st *sipTransportImpl, index TransportIndex) *UdpMessageSender {
	return &UdpMessageSender{
		sipTransport: st,
		index:        index,
		encoder:      enc.NewEncoder(),
	}
}

func (u *UdpMessageSender) SendRequest(request *sipmsg.GenericMessage) error {
	encodedMsg, err := u.encoder.Encode(request)
	if err != nil {
		return err
	}

	addr := &net.UDPAddr{
		IP:   net.ParseIP(u.index.Ip),
		Port: u.index.Port,
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}

	u.sipTransport.AddConn(u.index, conn)

	_, err = conn.Write(encodedMsg)
	return err
}

func (u *UdpMessageSender) SendResponse(response *sipmsg.GenericMessage) error {
	//TODO implement me
	panic("implement me")
}
