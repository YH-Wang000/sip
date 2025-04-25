package sipmsg

import (
	"bufio"
	"bytes"
)

type GenericMessage struct {
	StartLine     SipStartLine
	MessageHeader SipMessageHeader
	MessageBody   []byte
}

func ReadBody(s *bufio.Scanner, contentLength int) ([]byte, error) {
	buff := make([]byte, 0, contentLength)
	body := bytes.NewBuffer(buff)
	for s.Scan() {
		content := s.Text()
		if body.Len()+len(content) >= contentLength {
			body.WriteString(content[:contentLength-body.Len()])
			break
		}
		body.WriteString(content)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return body.Bytes(), nil
}

func IsRequest(message *GenericMessage) bool {
	return message.StartLine.IsRequestLine()
}

func IsResponse(message *GenericMessage) bool {
	return message.StartLine.IsStatusLine()
}
