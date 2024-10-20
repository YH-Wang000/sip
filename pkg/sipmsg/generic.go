package sipmsg

import "bufio"

type GenericMessage struct {
	StartLine     SipStartLine
	MessageHeader SipMessageHeader
	MessageBody   []byte
}

func ReadBody(s *bufio.Scanner, contentLength int) ([]byte, error) {
	return nil, nil
}
