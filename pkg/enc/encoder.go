package enc

import (
	"bytes"
	"sort"
	"strconv"

	"sip/pkg/sipmsg"
)

type Encoder interface {
	Encode(data *sipmsg.GenericMessage) ([]byte, error)
	SetHeaderOrderFunc(func(string, string) bool)
}

type EncoderOption func(e Encoder)

func NewEncoder(opts ...EncoderOption) Encoder {
	e := &encoder{}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

var DefaultHeadersSortOrderFunc = func(hKey1 string, hKey2 string) bool {
	return hKey1 < hKey2
}

func WithHeadersOrder(headersOrderFunc func(string, string) bool) EncoderOption {
	return func(e Encoder) {
		e.SetHeaderOrderFunc(headersOrderFunc)
	}
}

type encoder struct {
	// It is best to use it only in debug mode. Sorting requires additional performance.
	headersOrderFunc func(string, string) bool
}

func (e *encoder) Encode(msg *sipmsg.GenericMessage) ([]byte, error) {
	data := bytes.NewBuffer(nil)
	data.WriteString(msg.StartLine.String() + sipmsg.CRLF)
	// encoder attach content-length header
	if _, ok := msg.MessageHeader.Lookup("Content-Length"); !ok {
		msg.MessageHeader["Content-Length"] = []*sipmsg.HeaderFiledValue{
			{
				FieldValue: []string{strconv.Itoa(len(msg.MessageBody))},
				Params:     make(map[string]string),
			},
		}
	}
	if e.headersOrderFunc == nil {
		for headerName, headerValues := range msg.MessageHeader {
			for _, headerValue := range headerValues {
				data.WriteString(headerName + ": ")
				for i, value := range headerValue.FieldValue {
					data.WriteString(value)
					if i != len(headerValue.FieldValue)-1 {
						data.WriteString(",")
					}
				}
				for key, value := range headerValue.Params {
					data.WriteString(";" + key + "=" + value)
				}
				data.WriteString(sipmsg.CRLF)
			}
		}
	} else {
		sortedEncodeHeaders(data, msg.MessageHeader, e.headersOrderFunc)
	}
	data.WriteString(sipmsg.CRLF)
	data.Write(msg.MessageBody)
	return data.Bytes(), nil
}

func sortedEncodeHeaders(data *bytes.Buffer, headers sipmsg.SipMessageHeader, orderFunc func(string, string) bool) {
	headerKeys := make([]string, 0, len(headers))
	for headerKey := range headers {
		headerKeys = append(headerKeys, headerKey)
	}
	// todo Golang 1.22 version can be replaced with faster [slices.SortFunc]
	sort.Slice(headerKeys, func(i, j int) bool {
		return orderFunc(headerKeys[i], headerKeys[j])
	})
	for _, headerKey := range headerKeys {
		for _, fieldValue := range headers[headerKey] {
			data.WriteString(headerKey + ": ")
			for i, value := range fieldValue.FieldValue {
				data.WriteString(value)
				if i != len(fieldValue.FieldValue)-1 {
					data.WriteString(",")
				}
			}
			for key, value := range fieldValue.Params {
				data.WriteString(";" + key + "=" + value)
			}
			data.WriteString(sipmsg.CRLF)
		}
	}
}

func (e *encoder) SetHeaderOrderFunc(f func(string, string) bool) {
	e.headersOrderFunc = f
}
