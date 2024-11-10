package sipmsg

const (
	CRLF = "\r\n"
	SP   = " "
	TAB  = "\t"
)

type MethodEnum string

const (
	Register  MethodEnum = "REGISTER"
	Invite    MethodEnum = "INVITE"
	Ack       MethodEnum = "ACK"
	Cancel    MethodEnum = "CANCEL"
	Bye       MethodEnum = "BYE"
	Options   MethodEnum = "OPTIONS"
	Subscribe MethodEnum = "SUBSCRIBE"
	Notify    MethodEnum = "NOTIFY"
)

var SipRequestMethodSet = map[MethodEnum]bool{
	Register:  true,
	Invite:    true,
	Ack:       true,
	Cancel:    true,
	Bye:       true,
	Options:   true,
	Subscribe: true,
	Notify:    true,
}

type SipSchemeEnum string

const (
	Sip  SipSchemeEnum = "sip"
	Sips SipSchemeEnum = "sips"
	Tel  SipSchemeEnum = "tel"
)

type SipUriParamNameEnum string

const (
	TransportKey SipUriParamNameEnum = "transport"
	MaddrKey     SipUriParamNameEnum = "maddr"
	TtlKey       SipUriParamNameEnum = "ttl"
	UserKey      SipUriParamNameEnum = "user"
	MethodKey    SipUriParamNameEnum = "method"
	LrKey        SipUriParamNameEnum = "lr"
)

const (
	DefaultSipVersion = "SIP/2.0"
	DefaultSipPort    = "5060"
	DefaultSipsPort   = "5060"
)
