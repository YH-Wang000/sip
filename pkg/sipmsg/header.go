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

	headerName, text, found := strings.Cut(text, ":")
	if !found {
		return "", nil, errors.New("invalid header string")
	}

	return headerName, paresHeaderFiledValue(text), nil
}

func paresHeaderFiledValue(text string) *HeaderFiledValue {
	text = strings.TrimSpace(text)
	indices := getAngleBracketIndices(text)
	result := &HeaderFiledValue{
		FiledValue: make([]string, 0),
		Params:     make(map[string]string),
	}
	state := 0
	lastIndex := 0
	for i, c := range text {
		if c == ',' && state == 0 && checkIndices(indices, i) {
			value := text[lastIndex:i]
			result.FiledValue = append(result.FiledValue, strings.TrimSpace(value))
			lastIndex = i + 1
		}
		if c == ';' && checkIndices(indices, i) {
			if state == 0 {
				value := text[lastIndex:i]
				result.FiledValue = append(result.FiledValue, value)
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
		result.FiledValue = append(result.FiledValue, strings.TrimSpace(text[lastIndex:]))
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

// getAngleBracketIndices 从字符串中获取所有成对的 < 和 > 的索引
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
