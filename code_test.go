package gogeo

import (
	"encoding/json"
	"errors"
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

func Test_Code2ISO2(t *testing.T) {
	var tests = []struct {
		cc   Code2
		want string
	}{
		{Code2{'U', 'S'}, "US"},
		{Code2{'U', 'K'}, "UK"},
		{Code2{'A', '1'}, "A1"},
		{UndefinedCountryCode2, UndefinedCountryCodeISO2},
		{Code2{'-', '-'}, UndefinedCountryCodeISO2},
		{Code2{'X', 'Y'}, UndefinedCountryCodeISO2},
		{Code2{'u', 's'}, UndefinedCountryCodeISO2},
	}

	for _, test := range tests {
		if got := test.cc.ISO2(); got != test.want {
			t.Errorf("Code2(%q).ISO2() = %q, want %q", string(test.cc[:]), got, test.want)
		}
		if got := test.cc.String(); got != test.want {
			t.Errorf("Code2(%q).String() = %q, want %q", string(test.cc[:]), got, test.want)
		}
	}
}

func Test_Code2Scan(t *testing.T) {
	var cc Code2
	if err := cc.Scan([]byte("US")); err != nil {
		t.Fatalf("scan bytes: %v", err)
	}
	if cc != (Code2{'U', 'S'}) {
		t.Fatalf("scan bytes got %q", cc)
	}
	if err := cc.Scan("DE"); err != nil {
		t.Fatalf("scan string: %v", err)
	}
	if cc != (Code2{'D', 'E'}) {
		t.Fatalf("scan string got %q", cc)
	}
	if err := cc.Scan([]byte("USA")); !errors.Is(err, ErrCode2InvalidValueSize) {
		t.Fatalf("scan long bytes err = %v", err)
	}
	if err := cc.Scan(""); !errors.Is(err, ErrCode2InvalidValueSize) {
		t.Fatalf("scan empty err = %v", err)
	}
	if err := cc.Scan(123); !errors.Is(err, ErrCode2InvalidScanType) {
		t.Fatalf("scan int err = %v", err)
	}
}
