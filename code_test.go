package gogeo

import (
	"encoding/json"
	"testing"
)

func Test_JSONEncodingDecoding(t *testing.T) {
	type object struct {
		Code2 Code2 `json:"code2"`
	}

	var tests = []struct {
		source string
		result string
	}{
		{
			source: `{"code2":"**"}`,
			result: `{"code2":"**"}`,
		},
		{
			source: `{"code2":"A1"}`,
			result: `{"code2":"A1"}`,
		},
		{
			source: `{"code2":"***"}`,
			result: `{"code2":"**"}`,
		},
		{
			source: `{"code2":"ABC"}`,
			result: `{"code2":"**"}`,
		},
	}

	for _, test := range tests {
		var (
			obj  object
			data []byte
			err  error
		)
		if err = json.Unmarshal([]byte(test.source), &obj); err != nil {
			t.Errorf("invalid object unmarshal: %s", err.Error())
		}
		if data, err = json.Marshal(&obj); err != nil {
			t.Errorf("invalid object marshal: %s", err.Error())
		}

		if string(data) != test.result {
			t.Errorf(`invalid object conversion [%s] must be [%s]`, string(data), test.result)
		}
	}
}

func Test_CountryCode2ByString(t *testing.T) {
	var tests = []struct {
		code        string
		resultCode3 string
	}{
		{
			code:        "**",
			resultCode3: "***",
		},
		{
			code:        "--",
			resultCode3: "***",
		},
		{
			code:        "X",
			resultCode3: "***",
		},
		{
			code:        "UK",
			resultCode3: "GBR",
		},
		{
			code:        "GGY",
			resultCode3: "GGY",
		},
	}

	for _, test := range tests {
		if cc := CountryCode2ByString(test.code); cc.ISO3() != test.resultCode3 {
			t.Errorf(`invalid country code [%s] must be [%s]`, cc.ISO3(), test.resultCode3)
		}
	}
}
