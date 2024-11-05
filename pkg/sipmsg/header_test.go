package sipmsg

import (
	"reflect"
	"testing"
)

type headerTestCaseArgs struct {
	text string
}

type headerTestCase struct {
	name            string
	args            headerTestCaseArgs
	wantHeaderKey   string
	wantHeaderValue *HeaderFiledValue
	wantErr         bool
}

func TestParseHeader(t *testing.T) {

	testCases := getHeaderTestCases()

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ParseHeader(tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantHeaderKey {
				t.Errorf("ParseHeader() got = %v, want %v", got, tt.wantHeaderKey)
			}
			if !reflect.DeepEqual(got1, tt.wantHeaderValue) {
				t.Errorf("ParseHeader() got1 = %v, want %v", got1, tt.wantHeaderValue)
			}
		})
	}
}

func getHeaderTestCases() []headerTestCase {
	cases := make([]headerTestCase, 0)
	cases = append(cases, getFromHeaderTestCase())
	cases = append(cases, getToHeaderTestCase())
	cases = append(cases, getCSeqHeaderTestCase())
	cases = append(cases, getViaHeaderTestCase())
	cases = append(cases, getMaxForwardsHeaderTestCase())
	cases = append(cases, getUserAgentHeaderTestCase())
	cases = append(cases, getExpiresHeaderTestCase())
	cases = append(cases, getContactHeaderTestCase())
	return cases
}

func getFromHeaderTestCase() headerTestCase {
	return headerTestCase{
		name: "parse header From",
		args: headerTestCaseArgs{
			text: "From: \"123456\" <sip:123456@127.0.0.1>;tag=24d8f9d0",
		},
		wantHeaderKey: "From",
		wantHeaderValue: &HeaderFiledValue{
			FiledValue: []string{"\"123456\" <sip:123456@127.0.0.1>"},
			Params: map[string]string{
				"tag": "24d8f9d0",
			},
		},
		wantErr: false,
	}
}

func getToHeaderTestCase() headerTestCase {
	return headerTestCase{
		name: "parse header To",
		args: headerTestCaseArgs{
			text: "To: \"123456\" <sip:123456@127.0.0.1>",
		},
		wantHeaderKey: "To",
		wantHeaderValue: &HeaderFiledValue{
			FiledValue: []string{"\"123456\" <sip:123456@127.0.0.1>"},
			Params:     map[string]string{},
		},
		wantErr: false,
	}
}

func getCSeqHeaderTestCase() headerTestCase {
	return headerTestCase{
		name: "parse header CSeq",
		args: headerTestCaseArgs{
			text: "CSeq: 1 REGISTER",
		},
		wantHeaderKey: "CSeq",
		wantHeaderValue: &HeaderFiledValue{
			FiledValue: []string{"1 REGISTER"},
			Params:     map[string]string{},
		},
		wantErr: false,
	}
}

func getViaHeaderTestCase() headerTestCase {
	return headerTestCase{
		name: "parse header Via",
		args: headerTestCaseArgs{
			text: "Via: SIP/2.0/UDP 127.0.0.1:5060;branch=z9hG4bK-363233-d3bac43a3ab7f10841235227ef285e89",
		},
		wantHeaderKey: "Via",
		wantHeaderValue: &HeaderFiledValue{
			FiledValue: []string{"SIP/2.0/UDP 127.0.0.1:5060"},
			Params: map[string]string{
				"branch": "z9hG4bK-363233-d3bac43a3ab7f10841235227ef285e89",
			},
		},
		wantErr: false,
	}
}

func getMaxForwardsHeaderTestCase() headerTestCase {
	return headerTestCase{
		name: "parse header Max-Forwards",
		args: headerTestCaseArgs{
			text: "Max-Forwards: 70",
		},
		wantHeaderKey: "Max-Forwards",
		wantHeaderValue: &HeaderFiledValue{
			FiledValue: []string{"70"},
			Params:     map[string]string{},
		},
		wantErr: false,
	}
}

func getUserAgentHeaderTestCase() headerTestCase {
	return headerTestCase{
		name: "parse header User-Agent",
		args: headerTestCaseArgs{
			text: "User-Agent: Jitsi2.10.5550Windows 10",
		},
		wantHeaderKey: "User-Agent",
		wantHeaderValue: &HeaderFiledValue{
			FiledValue: []string{"Jitsi2.10.5550Windows 10"},
			Params:     map[string]string{},
		},
		wantErr: false,
	}
}

func getExpiresHeaderTestCase() headerTestCase {
	return headerTestCase{
		name: "parse header Expires",
		args: headerTestCaseArgs{
			text: "Expires: 600",
		},
		wantHeaderKey: "Expires",
		wantHeaderValue: &HeaderFiledValue{
			FiledValue: []string{"600"},
			Params:     map[string]string{},
		},
		wantErr: false,
	}
}

func getContactHeaderTestCase() headerTestCase {
	return headerTestCase{
		name: "parse header Contact",
		args: headerTestCaseArgs{
			text: "Contact: \"123456\" <sip:123456@127.0.0.1:5060;transport=udp;registering_acc=127_0_0_1>;expires=600",
		},
		wantHeaderKey: "Contact",
		wantHeaderValue: &HeaderFiledValue{
			FiledValue: []string{"\"123456\" <sip:123456@127.0.0.1:5060;transport=udp;registering_acc=127_0_0_1>"},
			Params: map[string]string{
				"expires": "600",
			},
		},
		wantErr: false,
	}
}
