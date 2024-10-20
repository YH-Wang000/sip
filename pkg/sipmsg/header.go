package sipmsg

import (
	"bufio"
	"errors"
	"maps"
	"strings"
)

// SipMessageHeader header = "header-name" HCOLON header-value *(COMMA header-value)
// HCOLON = ':'+' ' = ": "
// header-name: header-value1,header-value2;param1=value1;param2=value2
type SipMessageHeader map[string]*HeaderFiledValue

type HeaderFiledValue struct {
	FiledValue []string
	Params     map[string]string
}

func (h SipMessageHeader) Lookup(key string) (*HeaderFiledValue, bool) {
	value, ok := h[key]
	return value, ok
}

func ParseHeader(text string) (string, *HeaderFiledValue, error) {
	if len(text) == 0 {
		return "", nil, errors.New("invalid header string")
	}

	colonIndex := strings.IndexByte(text, ':')
	if colonIndex < 0 {
		return "", nil, errors.New("invalid header string")
	}
	headerName := strings.TrimSpace(text[:colonIndex])
	text = strings.TrimSpace(text[colonIndex+1:])
	semicolonIndex := strings.IndexByte(text, ';')
	if semicolonIndex < 0 {
		return headerName, &HeaderFiledValue{
			FiledValue: strings.Split(text, ","),
		}, nil
	}
	filedValueStr := text[:semicolonIndex]
	paramsStr := text[semicolonIndex+1:]

	headerFields := strings.Split(filedValueStr, ",")

	params := make(map[string]string)
	paramsSplit := strings.Split(paramsStr, ";")
	for i := 0; i < len(paramsSplit); i++ {
		kvSplit := strings.Split(paramsSplit[i], "=")
		if len(kvSplit) == 2 {
			params[kvSplit[0]] = kvSplit[1]
		} else {
			params[kvSplit[0]] = ""
		}
	}
	return headerName, &HeaderFiledValue{
		FiledValue: headerFields,
		Params:     params,
	}, nil
}

func ReadHeaders(s *bufio.Scanner) (map[string]*HeaderFiledValue, error) {
	headers := make(map[string]*HeaderFiledValue)
	var wholeHeader strings.Builder
	for s.Scan() {
		headerLine := s.Text()
		if headerLine == CRLF {
			// This means the end of the headers
			break
		}
		if !strings.HasPrefix(headerLine, SP) && !strings.HasPrefix(headerLine, TAB) {
			// This means a new header, the old header must be processed
			if whs := wholeHeader.String(); len(whs) > 0 {
				headerKey, headerValue, err := ParseHeader(whs)
				if err != nil {
					return nil, err
				}
				if hValue, ok := headers[headerKey]; ok {
					headerValue.FiledValue = append(headerValue.FiledValue, hValue.FiledValue...)
					maps.Copy(headerValue.Params, hValue.Params)
				}
				headers[headerKey] = headerValue
			}
			// Clear the previous row of data and prepare the next header
			wholeHeader.Reset()
		}
		headerLine = strings.TrimPrefix(headerLine, SP)
		headerLine = strings.TrimPrefix(headerLine, TAB)
		wholeHeader.WriteString(headerLine)
	}
	if s.Err() != nil {
		return nil, s.Err()
	}
	return headers, nil
}
