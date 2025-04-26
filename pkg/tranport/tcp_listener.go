package tranport

import (
	"context"
	"net"
	"sync"

	"sip/pkg/enc"
	"sip/pkg/log"
	"sip/pkg/sipmsg"
	"sip/pkg/util"
)

type tcpSocketListener struct {
	conn       net.Conn
	connStatus bool
	connMutex  sync.Mutex
	sourceAddr string
	ctx        context.Context
	cancel     context.CancelFunc

	msgC     chan *sipmsg.GenericMessage
	listener net.Listener
}

func NewTcpSocketListener(conn net.Conn) SocketListener {
	tl := &tcpSocketListener{
		conn:       conn,
		connStatus: true,
		sourceAddr: conn.LocalAddr().String(),
		msgC:       make(chan *sipmsg.GenericMessage, 1024),
	}
	tl.ctx, tl.cancel = context.WithCancel(context.Background())
	return tl
}

func (t *tcpSocketListener) ReadMessage() (*sipmsg.GenericMessage, error) {
	for {
		select {
		case <-t.ctx.Done():
			log.Debug("tcp socket listener closed")
			return nil, t.ctx.Err()
		case msg := <-t.msgC:
			return msg, nil
		}
	}
}

func (t *tcpSocketListener) Start() {
	go t.readLoop(t.conn)
	go t.listenLoop()
}

func (t *tcpSocketListener) Close() {
	t.cancel()
	util.Silently(t.listener.Close())
}

func (t *tcpSocketListener) getConnStatus() bool {
	t.connMutex.Lock()
	defer t.connMutex.Unlock()
	return t.connStatus
}

func (t *tcpSocketListener) setConnStatus(isValid bool) {
	t.connMutex.Lock()
	defer t.connMutex.Unlock()
	t.connStatus = isValid
}

func (t *tcpSocketListener) readLoop(conn net.Conn) {
	defer conn.Close()
	// todo set conn read deadline
	for {
		select {
		case <-t.ctx.Done():
			log.Debug("tcp socket listener closed")
			return
		default:
		}
		decoder := enc.NewDecoder(conn)
		message, err := decoder.ReadMessage() // if no conn deadline is set, maybe blocked forever
		if err != nil {
			log.Error("read message error: %v", err)
			t.setConnStatus(false)
			return
		}
		log.Debug("read message: %v", message)
		t.msgC <- message
	}
}

func (t *tcpSocketListener) listenLoop() {
	listener, err := net.Listen("tcp", t.sourceAddr)
	if err != nil {
		log.Error("listen error: %v", err)
		t.cancel()
		return
	}
	t.listener = listener
	defer listener.Close()
	for {
		select {
		case <-t.ctx.Done():
			log.Debug("tcp socket listener closed")
			return
		default:
		}
		conn, errAccept := listener.Accept()
		if errAccept != nil {
			log.Error("accept error: %v", errAccept)
			t.cancel()
			return
		}

		if t.getConnStatus() {
			util.Silently(conn.Close())
			continue
		}
		log.Debug("accept new conn: %v", conn.RemoteAddr())
		t.conn = conn
		go t.readLoop(t.conn)
	}
}
