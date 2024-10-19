package enc

import (
	"bufio"
	"bytes"
	"net"
	"strings"
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
	// initial capacity 1KB, max capacity 8KB
	if d.bodyLengthLimit+d.headerLengthLimit == 0 {
		d.SetBodyLengthLimit(DefaultMessageLength / 2)
		d.SetBodyLengthLimit(DefaultMessageLength / 2)
	}

	d.s.Buffer(make([]byte, InitDecoderBufferSize), d.headerLengthLimit+d.bodyLengthLimit)
	d.s.Split(splitCRLF)
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
	d.readHeaderTimeout = timeout
}

func (d *decoder) SetReadBodyTimeout(timeout time.Duration) {
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
	// read headers
	headers := make(map[string][]string)
	var wholeHeader strings.Builder
	for d.s.Scan() {
		headerLine := d.s.Text()
		if headerLine == sipmsg.CRLF {
			// This means the end of the headers
			break
		}
		if !strings.HasPrefix(headerLine, sipmsg.SP) && !strings.HasPrefix(headerLine, sipmsg.TAB) {
			// This means a new header, the old header must be processed
			headerKey, headerValue, err := sipmsg.ParseHeader(wholeHeader.String())
			if err != nil {
				return nil, err
			}
			if hValue, ok := headers[headerKey]; ok {
				// Need to merge repeated header key
				headerValue = append(headerValue, hValue...)
			}
			headers[headerKey] = headerValue
			// Clear the previous row of data and prepare the next header
			wholeHeader.Reset()
		}
		wholeHeader.WriteString(headerLine)
	}
	msg.MessageHeader = headers

	if noBody(msg.StartLine) {
		return msg, nil
	}

	err = d.conn.SetReadDeadline(time.Now().Add(d.readBodyTimeout))
	if err != nil {
		return nil, err
	}
	// read body
	// todo just read according to the Content-Length header

	return nil, nil
}

func noBody(startLine sipmsg.SipStartLine) bool {
	// todo Check all cases where body is not required
	if startLine.IsRequestLine() {
		return startLine.(*sipmsg.RequestLine).Method == sipmsg.Ack
	}
	if startLine.IsStatusLine() {
		return startLine.(*sipmsg.StatusLine).StatusCode < 200
	}
	return true
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
