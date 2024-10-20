package enc

import (
	"bufio"
	"bytes"
	"errors"
	"net"
	"strconv"
	"time"

	"sip/pkg/sipmsg"
)

type Decoder interface {
	ReadMessage() (*sipmsg.GenericMessage, error)
	SetHeaderLengthLimit(limit int)
	SetBodyLengthLimit(limit int)
	SetReadHeaderTimeout(limit time.Duration)
	SetReadBodyTimeout(limit time.Duration)
}

type DecoderOption func(d Decoder)

func NewDecoder(conn net.Conn, opts ...DecoderOption) Decoder {
	d := &decoder{
		conn: conn,
		s:    bufio.NewScanner(conn),
	}
	for _, opt := range opts {
		opt(d)
	}
	if d.bodyLengthLimit <= 0 {
		d.SetBodyLengthLimit(DefaultMessageLength / 2)
	}
	if d.headerLengthLimit <= 0 {
		d.SetHeaderLengthLimit(DefaultMessageLength / 2)
	}

	d.s.Buffer(make([]byte, InitDecoderBufferSize), d.headerLengthLimit+d.bodyLengthLimit)
	d.s.Split(splitCRLF)

	if d.readHeaderTimeout <= 0 {
		d.SetReadHeaderTimeout(DefaultReadHeaderTimeout)
	}
	if d.readBodyTimeout <= 0 {
		d.SetReadHeaderTimeout(DefaultReadBodyTimeout)
	}
	return d
}

func WithHeaderLengthLimit(limit int) DecoderOption {
	return func(d Decoder) {
		d.SetHeaderLengthLimit(limit)
	}
}
func WithBodyLengthLimit(limit int) DecoderOption {
	return func(d Decoder) {
		d.SetBodyLengthLimit(limit)
	}
}

func WithReadHeaderTimeout(timeout time.Duration) DecoderOption {
	return func(d Decoder) {
		d.SetReadHeaderTimeout(timeout)
	}
}
func WithReadBodyTimeout(timeout time.Duration) DecoderOption {
	return func(d Decoder) {
		d.SetReadBodyTimeout(timeout)
	}
}

type decoder struct {
	conn              net.Conn
	s                 *bufio.Scanner
	readHeaderTimeout time.Duration
	readBodyTimeout   time.Duration
	headerLengthLimit int
	bodyLengthLimit   int
}

func (d *decoder) SetHeaderLengthLimit(limit int) {
	d.headerLengthLimit = limit
}

func (d *decoder) SetBodyLengthLimit(limit int) {
	d.bodyLengthLimit = limit
}

func (d *decoder) SetReadHeaderTimeout(timeout time.Duration) {
	if d.readHeaderTimeout <= 0 {
		d.readHeaderTimeout = DefaultReadHeaderTimeout
		return
	}
	d.readHeaderTimeout = timeout
}

func (d *decoder) SetReadBodyTimeout(timeout time.Duration) {
	if d.readBodyTimeout <= 0 {
		d.readBodyTimeout = DefaultReadBodyTimeout
		return
	}
	d.readBodyTimeout = timeout
}

func (d *decoder) ReadMessage() (*sipmsg.GenericMessage, error) {
	msg := &sipmsg.GenericMessage{}
	if !d.s.Scan() {
		return nil, d.s.Err()
	}
	startLine, err := sipmsg.ParseStartLine(d.s.Text())
	if err != nil {
		return nil, err
	}
	msg.StartLine = startLine

	err = d.conn.SetReadDeadline(time.Now().Add(d.readHeaderTimeout))
	if err != nil {
		return nil, err
	}
	headers, err := sipmsg.ReadHeaders(d.s)
	if err != nil {
		return nil, err
	}
	msg.MessageHeader = headers

	if noBody(msg) {
		return msg, nil
	}
	err = d.conn.SetReadDeadline(time.Now().Add(d.readBodyTimeout))
	if err != nil {
		return nil, err
	}
	contentLength, err := getContentLength(msg.MessageHeader)
	body, err := sipmsg.ReadBody(d.s, contentLength)
	if err != nil {
		return nil, err
	}
	msg.MessageBody = body
	err = d.conn.SetReadDeadline(time.Time{})
	if err != nil {
		return nil, err
	}

	err = sipmsg.Validator.Struct(msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func getContentLength(header sipmsg.SipMessageHeader) (int, error) {
	headerValue, ok := header.Lookup("Content-Length")
	if !ok {
		return 0, errors.New("Content-Length header not found")
	}
	if len(headerValue.FiledValue) == 0 {
		return 0, errors.New("Content-Length header value is empty")
	}
	return strconv.Atoi(headerValue.FiledValue[0])
}

func noBody(msg *sipmsg.GenericMessage) bool {
	// todo Check all cases where body is not required
	contentLength, err := getContentLength(msg.MessageHeader)
	if err != nil {
		return true
	}
	if contentLength == 0 {
		return true
	}
	return false
}

func splitCRLF(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if i := bytes.Index(data, []byte("\r\n")); i >= 0 {
		return i + 2, data[:i], nil
	}

	if atEOF && len(data) > 0 {
		return len(data), data, nil
	}

	return 0, nil, nil
}
