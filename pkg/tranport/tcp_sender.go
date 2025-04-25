package tranport

import (
	"net"
	"strconv"
	"time"

	"sip/pkg/enc"
	"sip/pkg/log"
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
	// todo If the Role is a proxy, 5060 should be used as the source port first, for reuse conn between proxies

	// get connection
	u.sipTransport.ConnMapMutex.Lock()
	conn, ok := u.sipTransport.GetConn(u.index)
	if !ok {
		conn, err = u.createNewConn()
		if err != nil {
			return err
		}
		u.sipTransport.AddConn(u.index, conn)
	}
	u.sipTransport.ConnMapMutex.Unlock()

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

func (u *TcpMessageSender) createNewConn() (net.Conn, error) {
	if u.sipTransport.Role == "proxy" {
		dialer := &net.Dialer{
			LocalAddr: &net.TCPAddr{
				IP:   net.ParseIP("0.0.0.0"),
				Port: u.proxyDefaultSourcePort,
			},
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}
		conn, err := dialer.Dial(u.index.Network, u.index.Ip+":"+strconv.Itoa(u.index.Port))
		if err == nil {
			// proxy should use 5060 as source port first
			return conn, nil
		}
		log.Debug("5060 is already used, proxy can't use it")
	}
	return net.Dial(u.index.Network, u.index.Ip+":"+strconv.Itoa(u.index.Port))
}
