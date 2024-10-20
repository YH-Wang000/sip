package sipmsg

import (
	"errors"
	"strings"
)

type SipStartLine interface {
	String() string
	IsRequestLine() bool
	IsStatusLine() bool
}

func ParseStartLine(line string) (SipStartLine, error) {
	if len(line) == 0 {
		return nil, errors.New("start line is invalid")
	}
	if strings.HasPrefix(line, DefaultSipVersion) {

		return &StatusLine{}, nil
	}
	split := strings.Split(line, SP)
	if len(split) != 3 {
		return nil, errors.New("start line is invalid")
	}
	rl := &RequestLine{}
	rl.Method = MethodEnum(split[0])
	sipUri, err := ParseSipUri(split[1])
	if err != nil {
		return nil, err
	}
	rl.RequestUri = sipUri
	rl.SipVersion = split[2]
	return rl, nil
}

type RequestLine struct {
	Method     MethodEnum `validate:"required, oneof=REGISTER INVITE ACK CANCEL BYE OPTIONS"`
	RequestUri *SipUri    `validate:"required"`
	SipVersion string     `validate:"required, checkSipVersion"`
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
	Scheme        SipSchemeEnum `validate:"required oneof=sip sips tel"`
	User          string        // resource identifier in host
	Password      string        // not recommended may cause security risks
	Host          string        `validate:"required"` // The host providing the SIP resource
	Port          string        `validate:"numeric"`
	UriParameters map[SipUriParamNameEnum]string
	Headers       map[string]string
}

func ParseSipUri(uri string) (*SipUri, error) {
	sipUrl := &SipUri{}
	colonIndex := strings.IndexByte(uri, ':')
	if colonIndex < 0 {
		return nil, errors.New("invalid sip uri")
	}
	sipUrl.Scheme = SipSchemeEnum(uri[:colonIndex])
	uri = uri[colonIndex+1:]

	atIndex := strings.IndexByte(uri, '@')
	if atIndex != -1 {
		// optional userinfo
		userInfo := uri[:atIndex]
		userInfoSplit := strings.Split(userInfo, ":")
		if len(userInfoSplit) == 1 {
			sipUrl.User = userInfoSplit[0]
		} else if len(userInfoSplit) == 2 {
			sipUrl.User = userInfoSplit[0]
			sipUrl.Password = userInfoSplit[1]
		} else {
			return nil, errors.New("invalid sip uri, userinfo should be include user and password")
		}
		uri = uri[atIndex+1:]
	}

	semicolonIndex := strings.IndexByte(uri, ';')
	if semicolonIndex < 0 {
		sipUrl.Host, sipUrl.Port = ParseDomain(uri)
		return sipUrl, nil
	}
	domain := uri[:semicolonIndex]
	sipUrl.Host, sipUrl.Port = ParseDomain(domain)
	uri = uri[semicolonIndex+1:]

	questionIndex := strings.IndexByte(uri, '?')
	if questionIndex < 0 {
		sipUrl.UriParameters = ParseUriParams(uri)
		return sipUrl, nil
	}
	sipUrl.UriParameters = ParseUriParams(uri[:questionIndex])
	uri = uri[questionIndex+1:]

	if len(uri) != 0 {
		sipUrl.Headers = ParseSipHeader(uri)
	}

	return sipUrl, nil
}

func ParseSipHeader(uri string) map[string]string {
	return nil
}

func ParseUriParams(uri string) map[SipUriParamNameEnum]string {
	return nil
}

func ParseDomain(domain string) (string, string) {
	if len(domain) == 0 {
		return "", ""
	}
	split := strings.Split(domain, ":")
	if len(split) == 1 {
		return split[0], DefaultSipPort
	}
	return split[0], split[1]
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
