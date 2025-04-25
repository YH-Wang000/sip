package sipmsg

import "strings"

var hdrKeyMap = map[string]string{
	"accept":                      "a",
	"accept-encoding":             "ace",
	"accept-language":             "acp",
	"alert-info":                  "ali",
	"allow":                       "alw",
	"allow-events":                "ale",
	"authentication-info":         "aut",
	"authorization":               "auh",
	"call-id":                     "i",
	"call-info":                   "cal",
	"contact":                     "m",
	"content-disposition":         "dis",
	"content-encoding":            "e",
	"content-language":            "l",
	"content-length":              "len",
	"content-type":                "c",
	"cseq":                        "q",
	"date":                        "d",
	"error-info":                  "err",
	"event":                       "o",
	"expires":                     "exp",
	"from":                        "f",
	"history-info":                "hi",
	"in-reply-to":                 "irt",
	"join":                        "j",
	"max-forwards":                "x",
	"mime-version":                "v",
	"min-se":                      "min",
	"organization":                "org",
	"priority":                    "p",
	"proxy-authenticate":          "ppa",
	"proxy-authorization":         "ppau",
	"proxy-require":               "rq",
	"reason":                      "r",
	"record-route":                "rr",
	"refer-to":                    "ref",
	"referred-by":                 "rb",
	"reject-contact":              "rec",
	"remote-party-id":             "rpid",
	"reply-to":                    "rep",
	"require":                     "req",
	"retry-after":                 "ra",
	"route":                       "rte",
	"rseq":                        "rs",
	"security-client":             "sec",
	"security-server":             "ses",
	"security-verify":             "sev",
	"server":                      "s",
	"service-route":               "srv",
	"session-expires":             "exp",
	"sip-etag":                    "set",
	"sip-if-match":                "sim",
	"subject":                     "sub",
	"supported":                   "k",
	"timestamp":                   "ts",
	"to":                          "t",
	"trusted":                     "tr",
	"unsupported":                 "uns",
	"user-agent":                  "u",
	"via":                         "v",
	"warning":                     "w",
	"www-authenticate":            "waa",
	"p-asserted-identity":         "pai",
	"p-associated-uri":            "par",
	"p-called-party-id":           "pcp",
	"p-charging-function-address": "pcc",
	"p-charging-vector":           "pcv",
	"p-dcs-trace-party-id":        "ptp",
	"p-dcs-osps":                  "pdcs",
	"p-dcs-billing-info":          "pbi",
	"p-dcs-redirect":              "pdr",
	"p-dcs-redirect-context":      "pdc",
	"p-dcs-3gpp":                  "p3g",
	"p-dcs-3gpp2":                 "p3g2",
	"p-media-authorization":       "pma",
	"p-preferred-identity":        "ppi",
	"p-visited-network-id":        "pvn",
	"path":                        "pa",
	"privacy":                     "pri",
}

type HeaderKey string

func (h HeaderKey) Normalize() string {
	for fullKey, shortKey := range hdrKeyMap {
		strings.EqualFold(string(h), shortKey)
		if string(h) == shortKey {
			return fullKey
		}
	}
	// unknown key, no normalization
	return string(h)
}

func (h HeaderKey) Shorten() string {
	for fullKey, shortKey := range hdrKeyMap {
		if strings.EqualFold(string(h), fullKey) {
			return shortKey
		}
	}
	// unknown key, no shortening
	return string(h)
}
