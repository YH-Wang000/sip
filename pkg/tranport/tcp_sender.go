package tranport

import (
	"errors"

	"sip/pkg/enc"
	"sip/pkg/sipmsg"
)

type TcpMessageSender struct {
	sipTransport           *sipTransportImpl
	index                  TransportIndex
	encoder                enc.Encoder
	proxyDefaultSourcePort int
}

func NewTcpMessageSender(st *sipTransportImpl, index TransportIndex) *TcpMessageSender {
	return &TcpMessageSender{
		sipTransport:           st,
		index:                  index,
		encoder:                enc.NewEncoder(),
		proxyDefaultSourcePort: 5060,
	}
}

func (u *TcpMessageSender) SendRequest(request *sipmsg.GenericMessage) error {
	reqBytes, err := u.encoder.Encode(request)
	if err != nil {
		return err
	}

	// get connection
	conn, ok := u.sipTransport.GetConn(u.index)
	if !ok {
		return errors.New("no available conn")
	}

	_, err = conn.Write(reqBytes)
	if err != nil {
		return err
	}
	return nil
}

func (u *TcpMessageSender) SendResponse(response *sipmsg.GenericMessage) error {
	//TODO implement me
	panic("implement me")
}
