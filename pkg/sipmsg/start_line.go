package sipmsg

type SipStartLine interface {
	String() string
	IsRequestLine() bool
	IsStatusLine() bool
}

func ParseStartLine(line string) (SipStartLine, error) {
	return nil, nil
}

type RequestLine struct {
	Method     MethodEnum
	RequestUri string
	SipVersion string
}

func (r *RequestLine) String() string {
	return ""
}

func (r *RequestLine) IsRequestLine() bool {
	return true
}

func (r *RequestLine) IsStatusLine() bool {
	return false
}

type SipUri struct {
	Scheme        SipSchemeEnum
	User          string // resource identifier in host
	Password      string // not recommended may cause security risks
	Host          string // The host providing the SIP resource
	Port          string
	UriParameters map[SipUriParamNameEnum]string
	Headers       map[string]string
}

type StatusLine struct {
	SipVersion   string
	StatusCode   int
	ReasonPhrase string // Can be served based on the language specified by the Accept-Language header
}

func (s *StatusLine) String() string {
	return ""
}

func (s *StatusLine) IsRequestLine() bool {
	return false
}

func (s *StatusLine) IsStatusLine() bool {
	return true
}
