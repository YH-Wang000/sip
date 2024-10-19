package sipmsg

type GenericMessage struct {
	StartLine     SipStartLine
	MessageHeader SipMessageHeader
	MessageBody   string
}
