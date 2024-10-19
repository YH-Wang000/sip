package sipmsg

// SipMessageHeader header = "header-name" HCOLON header-value *(COMMA header-value)
// HCOLON = ':'+' ' = ": "
type SipMessageHeader map[string][]string

func ParseHeader(text string) (string, []string, error) {
	return "", nil, nil
}
