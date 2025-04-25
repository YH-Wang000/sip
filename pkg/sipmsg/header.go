package sipmsg

import (
	"strings"

	"sip/pkg/log"
)

// SipMessageHeader header = "header-name" HCOLON header-value *(COMMA header-value)
// HCOLON = ':'+' ' = ": "
// header-name: header-value1,header-value2;param1=value1;param2=value2
type SipMessageHeader map[string][]*HeaderFiledValue

type HeaderFiledValue struct {
	FieldValue []string
	Params     map[string]string
}

func (h SipMessageHeader) Lookup(key string) ([]*HeaderFiledValue, bool) {
	value, ok := h[key]
	return value, ok
}

func (h SipMessageHeader) Add(key string, value *HeaderFiledValue) {
	if _, ok := h[key]; !ok {
		h[key] = []*HeaderFiledValue{value}
		return
	}
	h[key] = append(h[key], value)
}

type ViaHeader struct {
	Protocol  string            // 协议版本，例如 "SIP/2.0"
	Transport string            // 传输协议，例如 "UDP", "TCP", "TLS"
	SentBy    string            // 发送方地址，例如 "192.168.1.1:5060"
	Branch    string            // 分支标识，例如 "z9hG4bK123456"
	Params    map[string]string // 其他参数
}

func (v *ViaHeader) From(headerValue *HeaderFiledValue) {
	if len(headerValue.FieldValue) == 0 {
		return
	}

	// 解析 FieldValue，格式为 "SIP/2.0/UDP 192.168.1.1:5060"
	parts := strings.Fields(headerValue.FieldValue[0])
	if len(parts) < 2 {
		return
	}

	// 解析协议版本和传输协议，例如 "SIP/2.0/UDP"
	protocolParts := strings.Split(parts[0], "/")
	if len(protocolParts) == 3 {
		v.Protocol = protocolParts[0] + "/" + protocolParts[1]
		v.Transport = protocolParts[2]
	} else {
		log.Error("Invalid via header format")
	}

	// 解析发送方地址，例如 "192.168.1.1:5060"
	v.SentBy = parts[1]

	// 解析参数，例如 branch=z9hG4bK123456
	v.Params = make(map[string]string)
	for key, value := range headerValue.Params {
		v.Params[key] = value
		if key == "branch" {
			v.Branch = value
		}
	}
}
