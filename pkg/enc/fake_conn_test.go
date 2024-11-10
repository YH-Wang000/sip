package enc

import (
	"bytes"
	"net"
	"time"
)

type fakeConn struct {
	buff *bytes.Buffer
}

func (f *fakeConn) Read(b []byte) (n int, err error) {
	return f.buff.Read(b)
}

func (f *fakeConn) Write(b []byte) (n int, err error) {
	return f.buff.Write(b)
}

func (f *fakeConn) Close() error {
	return nil
}

func (f *fakeConn) LocalAddr() net.Addr {
	return nil
}

func (f *fakeConn) RemoteAddr() net.Addr {
	return nil
}

func (f *fakeConn) SetDeadline(_ time.Time) error {
	return nil
}

func (f *fakeConn) SetReadDeadline(_ time.Time) error {
	return nil
}

func (f *fakeConn) SetWriteDeadline(_ time.Time) error {
	return nil
}
