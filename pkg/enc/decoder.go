package enc

import (
	"bufio"
	"bytes"
	"errors"
	"net"
	"strconv"
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
	headers, err := ReadHeaders(d.s)
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
	headerValues, ok := header.Lookup("Content-Length")
	if !ok || len(headerValues) == 0 {
		return 0, errors.New("Content-Length header not found")
	}
	if len(headerValues[0].FieldValue) == 0 {
		return 0, errors.New("Content-Length header value is empty")
	}
	return strconv.Atoi(headerValues[0].FieldValue[0])
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

func ParseHeader(text string) (string, *sipmsg.HeaderFiledValue, error) {
	if len(text) == 0 {
		return "", nil, errors.New("invalid header string")
	}

	headerName, text, found := strings.Cut(text, ":")
	if !found {
		return "", nil, errors.New("invalid header string")
	}

	return headerName, paresHeaderFiledValue(text), nil
}

func paresHeaderFiledValue(text string) *sipmsg.HeaderFiledValue {
	text = strings.TrimSpace(text)
	indices := getAngleBracketIndices(text)
	result := &sipmsg.HeaderFiledValue{
		FieldValue: make([]string, 0),
		Params:     make(map[string]string),
	}
	state := 0
	lastIndex := 0
	for i, c := range text {
		if c == ',' && state == 0 && checkIndices(indices, i) {
			value := text[lastIndex:i]
			result.FieldValue = append(result.FieldValue, strings.TrimSpace(value))
			lastIndex = i + 1
		}
		if c == ';' && checkIndices(indices, i) {
			if state == 0 {
				value := text[lastIndex:i]
				result.FieldValue = append(result.FieldValue, value)
				lastIndex = i + 1
				state = 1
				continue
			}
			paramsStr := text[lastIndex:i]
			key, value, found := strings.Cut(paramsStr, "=")
			if !found {
				continue
			}
			result.Params[key] = value
			lastIndex = i + 1
		}
	}
	if state == 0 {
		result.FieldValue = append(result.FieldValue, strings.TrimSpace(text[lastIndex:]))
	} else if state == 1 {
		paramsStr := text[lastIndex:]
		key, value, found := strings.Cut(paramsStr, "=")
		if !found {
			return result
		}
		result.Params[key] = value
	}
	return result
}

func checkIndices(indices [][2]int, i int) bool {
	for _, index := range indices {
		if index[0] <= i && i <= index[1] {
			return false
		}
	}
	return true
}

// getAngleBracketIndices Get the indices of all pairs of < and > from a string
func getAngleBracketIndices(s string) [][2]int {
	var indices [][2]int
	balance := 0
	left := 0
	for i, char := range s {
		if char == '<' {
			if balance == 0 {
				left = i
			}
			balance++
		} else if char == '>' {
			balance--
			if balance == 0 {
				// 找到一对匹配的 < 和 >
				indices = append(indices, [2]int{left, i})
			}
		}
	}
	return indices
}

func ReadHeaders(s *bufio.Scanner) (map[string][]*sipmsg.HeaderFiledValue, error) {
	headers := make(map[string][]*sipmsg.HeaderFiledValue)
	var wholeHeader strings.Builder
	for s.Scan() {
		headerLine := s.Text()
		if headerLine == "" {
			// This means the end of the headers
			if err := populateHeadersMap(wholeHeader.String(), headers); err != nil {
				return nil, err
			}
			break
		}
		if !strings.HasPrefix(headerLine, sipmsg.SP) && !strings.HasPrefix(headerLine, sipmsg.TAB) {
			// This means a new header, the old header must be processed
			if err := populateHeadersMap(wholeHeader.String(), headers); err != nil {
				return nil, err
			}
			// Clear the previous row of data and prepare the next header
			wholeHeader.Reset()
		}
		headerLine = strings.TrimPrefix(headerLine, sipmsg.SP)
		headerLine = strings.TrimPrefix(headerLine, sipmsg.TAB)
		wholeHeader.WriteString(headerLine)
	}
	if s.Err() != nil {
		return nil, s.Err()
	}
	return headers, nil
}

func populateHeadersMap(headerStr string, headers map[string][]*sipmsg.HeaderFiledValue) error {
	if len(headerStr) > 0 {
		headerKey, headerValue, err := ParseHeader(headerStr)
		if err != nil {
			return err
		}
		headerFiledValues, ok := headers[headerKey]
		if !ok {
			headerFiledValues = make([]*sipmsg.HeaderFiledValue, 0, 1)
		}
		headerFiledValues = append(headerFiledValues, headerValue)
		headers[headerKey] = headerFiledValues
	}
	return nil
}
