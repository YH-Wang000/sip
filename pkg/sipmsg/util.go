package sipmsg

func GetHeaderFieldValues(headers SipMessageHeader, headerName string) ([]*HeaderFiledValue, bool) {
	if headerFiledValue, ok := headers[HeaderKey(headerName).Normalize()]; ok {
		return headerFiledValue, true
	}
	if headerFiledValue, ok := headers[HeaderKey(headerName).Shorten()]; ok {
		return headerFiledValue, true
	}
	return nil, false
}

func GetFirstHeaderByKey(headers SipMessageHeader, key string) (*HeaderFiledValue, bool) {
	fieldValues, ok := GetHeaderFieldValues(headers, key)
	if !ok || len(fieldValues) == 0 {
		return nil, false
	}
	return fieldValues[0], true
}
