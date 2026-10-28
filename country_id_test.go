package gogeo

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestCountryIDByString(t *testing.T) {
	us := CountryByCode2("US")
	id := CountryIDByString("US")
	if id != CountryID(us.ID) || id.Object() != us || id.ISO2() != "US" || id.String() != "US" {
		t.Fatalf("US id=%d iso2=%q object=%v", id.ID(), id.ISO2(), id.Object())
	}
	if id.Code() != (CountryCode2{'U', 'S'}) {
		t.Fatalf("US Code() = %q", id.Code())
	}
	if CountryIDByString("USA") != id || id.ISO3() != "USA" {
		t.Fatalf("USA = %d iso3=%q", CountryIDByString("USA"), id.ISO3())
	}

	if CountryIDByString("UK") != CountryIDByString("GB") {
		t.Fatal("UK and GB must be the same CountryID")
	}
	if CountryIDByString("UK").Code() != (CountryCode2{'G', 'B'}) {
		t.Fatalf("UK Code() = %q, want GB", CountryIDByString("UK").Code())
	}

	for _, code := range []string{"", "X"} {
		got := CountryIDByString(code)
		if got.ID() != 0 || got.ISO2() != "**" || got.Object() == nil {
			t.Fatalf("CountryIDByString(%q) = %d %q %v", code, got.ID(), got.ISO2(), got.Object())
		}
	}

	past := CountryID(uint8(len(Countries)))
	if past.ID() != 0 || past.ISO2() != "**" || past.Object() == nil {
		t.Fatalf("out-of-range id = %d %q %v", past.ID(), past.ISO2(), past.Object())
	}
}

func TestCountryIDJSON(t *testing.T) {
	type object struct {
		ID CountryID `json:"id"`
	}

	var tests = []struct {
		source string
		result string
	}{
		{source: `{"id":"US"}`, result: `{"id":"US"}`},
		{source: `{"id":"**"}`, result: `{"id":"**"}`},
		{source: `{"id":"***"}`, result: `{"id":"**"}`},
		{source: `{"id":"ABC"}`, result: `{"id":"**"}`},
	}

	for _, test := range tests {
		var obj object
		if err := json.Unmarshal([]byte(test.source), &obj); err != nil {
			t.Fatalf("unmarshal %s: %v", test.source, err)
		}
		data, err := json.Marshal(&obj)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != test.result {
			t.Fatalf("json %s = %s, want %s", test.source, data, test.result)
		}
	}
}

func TestCountryIDScan(t *testing.T) {
	var id CountryID
	if err := id.Scan([]byte("DE")); err != nil {
		t.Fatalf("scan bytes: %v", err)
	}
	if id != CountryID(CountryByCode2("DE").ID) || id.ISO2() != "DE" {
		t.Fatalf("scan DE = %d %q", id.ID(), id.ISO2())
	}
	v, err := id.Value()
	if err != nil || v != "DE" {
		t.Fatalf("Value() = %v %v", v, err)
	}
	if err := id.Scan([]byte("USA")); !errors.Is(err, ErrCountryIDInvalidValueSize) {
		t.Fatalf("scan long bytes err = %v", err)
	}
	if err := id.Scan(123); !errors.Is(err, ErrCountryIDInvalidScanType) {
		t.Fatalf("scan int err = %v", err)
	}
}
