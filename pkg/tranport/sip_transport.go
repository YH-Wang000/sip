package tranport

import (
	"errors"
	"net"
	"strconv"
	"sync"
	"time"

	"sip/pkg/log"
	"sip/pkg/sipmsg"
)

type MessageSender interface {
	SendRequest(request *sipmsg.GenericMessage) error
	SendResponse(response *sipmsg.GenericMessage) error
}

type SocketListener interface {
	ReadMessage() (*sipmsg.GenericMessage, error)
	Start()
	Close()
}

type ConnMgr interface {
	GetConn(TransportIndex) (net.Conn, bool)
}

type SipTransport interface {
	ConnMgr
	SendMessage(msg *sipmsg.GenericMessage) error
	SetMessageHandler(func(msg *sipmsg.GenericMessage))
}

type sipTransportImpl struct {
	Role string

	connMapMutex sync.Mutex
	connMap      map[TransportIndex]net.Conn

	listenerMapMutex sync.Mutex
	listenerMap      map[TransportIndex]SocketListener
}

func NewSipTransportImpl(role string) SipTransport {
	return &sipTransportImpl{
		Role:    role,
		connMap: make(map[TransportIndex]net.Conn),
	}
}

func (s *sipTransportImpl) SendMessage(msg *sipmsg.GenericMessage) error {
	transportIndex, err := parseTransportIndex(msg)
	if err != nil {
		return err
	}
	msgSender := s.createMessageSender(transportIndex)
	if sipmsg.IsRequest(msg) {
		return msgSender.SendRequest(msg)
	}
	if sipmsg.IsResponse(msg) {
		return msgSender.SendResponse(msg)
	}
	return errors.New("unknown message type")
}

func (s *sipTransportImpl) createMessageSender(index TransportIndex) MessageSender {
	if index.Network == "tcp" {
		return NewTcpMessageSender(s, index)
	} else {
		return NewUdpMessageSender(s, index)
	}
}

func parseTransportIndex(msg *sipmsg.GenericMessage) (TransportIndex, error) {
	switch startLine := msg.StartLine.(type) {
	case *sipmsg.RequestLine:
		return GetTargetTransportIndexFromRequest(startLine, msg.MessageHeader)
	case *sipmsg.StatusLine:
		return GetTargetTransportIndexFromResponse(startLine, msg.MessageHeader)
	default:
		return TransportIndex{}, errors.New("unknown start line type")
	}
}

func (s *sipTransportImpl) SetMessageHandler(f func(msg *sipmsg.GenericMessage)) {
	//TODO implement me
	panic("implement me")
}

func (s *sipTransportImpl) GetConn(index TransportIndex) (net.Conn, bool) {
	s.connMapMutex.Lock()
	defer s.connMapMutex.Unlock()
	conn, ok := s.connMap[index]
	if !ok {
		newConn, err := s.createNewConn(index)
		if err != nil {
			log.Error("create new conn error: ", err)
			return nil, false
		}
		s.connMap[index] = newConn
		return newConn, true
	}
	return conn, true
}

func (s *sipTransportImpl) createNewConn(index TransportIndex) (net.Conn, error) {
	switch index.Network {
	case "tcp":
		return s.createNewTcpConn(index)
	case "udp":
	}
	return nil, nil
}

func (s *sipTransportImpl) createNewTcpConn(index TransportIndex) (net.Conn, error) {
	if s.Role == "proxy" {
		dialer := &net.Dialer{
			LocalAddr: &net.TCPAddr{
				IP:   net.ParseIP("0.0.0.0"),
				Port: 5060,
			},
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}
		conn, err := dialer.Dial(index.Network, index.Ip+":"+strconv.Itoa(index.Port))
		if err == nil {
			// proxy should use 5060 as source port first
			log.Info("new tcp conn: ", conn.LocalAddr(), " -> ", conn.RemoteAddr())
			s.listenSocket(conn, index)
			return conn, nil
		}
		log.Debug("5060 is already used, proxy can't use it")
	}
	conn, err := net.Dial(index.Network, index.Ip+":"+strconv.Itoa(index.Port))
	if err != nil {
		return nil, err
	}
	log.Info("new tcp conn: ", conn.LocalAddr(), " -> ", conn.RemoteAddr())
	s.listenSocket(conn, index)
	return conn, nil
}

func (s *sipTransportImpl) GetConnByIndex(index TransportIndex) (net.Conn, bool) {
	s.connMapMutex.Lock()
	defer s.connMapMutex.Unlock()
	conn, ok := s.connMap[index]
	return conn, ok
}

func (s *sipTransportImpl) listenSocket(conn net.Conn, index TransportIndex) {
	s.listenerMapMutex.Lock()
	defer s.listenerMapMutex.Unlock()
	socketListener := s.createSocketListener(conn, index)
	s.listenerMap[index] = socketListener
	socketListener.Start()
}

func (s *sipTransportImpl) createSocketListener(conn net.Conn, index TransportIndex) SocketListener {
	switch index.Network {
	case "tcp":
		return NewTcpSocketListener(conn)
	case "udp":
		return NewUdpSocketListener(conn)
	default:
	}
	log.Error("unknown network type: ", index.Network)
	return nil
}

var defaultPortMap = map[string]int{
	"udp": 5060,
	"tcp": 5060,
	"tls": 5061,
}

func GetTargetTransportIndexFromResponse(_ *sipmsg.StatusLine, headers sipmsg.SipMessageHeader) (TransportIndex, error) {
	var index TransportIndex

	// 获取 Via 头字段
	viaHeaderFieldValue, ok := sipmsg.GetFirstHeaderByKey(headers, "Via")
	if !ok {
		return index, errors.New("no via header found")
	}

	// 解析 Via 头字段
	viaHeader := &sipmsg.ViaHeader{}
	viaHeader.From(viaHeaderFieldValue)
	index.Network = viaHeader.Transport

	// 解析 SentBy 地址
	host, portStr, err := net.SplitHostPort(viaHeader.SentBy)
	if err != nil {
		return index, err
	}

	// 解析 IP 地址
	if len(net.ParseIP(host)) != 0 {
		index.Ip = host
	} else {
		ips, err := net.LookupIP(host)
		if err != nil {
			return index, err
		}
		if len(ips) == 0 {
			return index, errors.New("no ip found")
		}
		index.Ip = ips[0].String()
	}

	// 解析端口
	defaultPort, ok := defaultPortMap[index.Network]
	if !ok {
		return index, errors.New("no default port found for network")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil || port == 0 {
		port = defaultPort
	}
	index.Port = port

	return index, nil
}

func GetTargetTransportIndexFromRequest(requestLine *sipmsg.RequestLine, hdrs sipmsg.SipMessageHeader) (TransportIndex, error) {
	var index TransportIndex

	if len(net.ParseIP(requestLine.RequestUri.Host)) != 0 {
		index.Ip = requestLine.RequestUri.Host
	} else {
		ips, err := net.LookupIP(requestLine.RequestUri.Host)
		if err != nil {
			return index, err
		}
		if len(ips) == 0 {
			return index, errors.New("no ip found")
		} else {
			index.Ip = ips[0].String()
		}
	}

	// parse network must before than port
	viaHeader := &sipmsg.ViaHeader{}
	viaHeaderFieldValue, ok := sipmsg.GetFirstHeaderByKey(hdrs, "Via")
	if !ok {
		return index, errors.New("no via header found")
	}
	viaHeader.From(viaHeaderFieldValue)
	index.Network = viaHeader.Transport

	defaultPort, ok := defaultPortMap[index.Network]
	if !ok {
		return index, errors.New("no default port found for network")
	}
	port, err := strconv.Atoi(requestLine.RequestUri.Port)
	if err != nil {
		return index, err
	}
	if port == 0 {
		port = defaultPort
	}

	return index, nil
}
