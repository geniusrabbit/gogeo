package gogeo

import (
	"math"
	"strings"
	"testing"
)

func TestCountryByCode2(t *testing.T) {
	us := CountryByCode2("US")
	if us == nil || us.ISO2() != "US" || us.ISO3() != "USA" || us.Name != "United States" {
		t.Fatalf("US: %+v iso2=%q iso3=%q", us, us.ISO2(), us.ISO3())
	}
	if us.Continent() != "NA" {
		t.Fatalf("US continent = %q", us.Continent())
	}
	if got := CountryByCode2Bytes(CountryCode2{'U', 'S'}); got != us {
		t.Fatalf("bytes lookup mismatch")
	}
	if got := CountryByID(us.ID); got != us {
		t.Fatalf("id lookup mismatch")
	}

	undef := CountryByCode2("--")
	if undef == nil || undef.ISO2() != UndefinedCountryCodeISO2 {
		t.Fatalf("unknown ISO2 must be undefined, got %+v", undef)
	}
	if CountryByCode2("X") != undef || CountryByCode2("") != undef {
		t.Fatalf("short/empty ISO2 must be undefined")
	}
}

func TestCountryByCode3(t *testing.T) {
	us := CountryByCode3("USA")
	if us == nil || us.ISO2() != "US" {
		t.Fatalf("USA: %+v", us)
	}
	undef := CountryByCode3("")
	if undef == nil || undef.ISO2() != UndefinedCountryCodeISO2 {
		t.Fatalf("empty Code3 must be undefined, got %+v", undef)
	}
	if CountryByCode3("???") != undef || CountryByCode3("AA1").ISO2() != "A1" {
		t.Fatalf("special/unknown Code3 mismatch")
	}

	// Empty Code3 entries must not collide on "".
	if CountryByCode2("A2").ISO3() != "" || CountryByCode2("AP").ISO3() != "" {
		t.Fatalf("pseudo-countries must keep empty Code3")
	}
	if CountryByCode3("") == CountryByCode2("A2") || CountryByCode3("") == CountryByCode2("EU") {
		t.Fatalf("empty Code3 must not resolve to a pseudo-country")
	}
}

func TestCountryPackedLists(t *testing.T) {
	us := CountryByCode2("US")
	if len(us.Currency()) < 1 || us.Currency()[0] != "USD" {
		t.Fatalf("US currency = %v", us.Currency())
	}
	if len(us.Languages()) != 1 || us.Languages()[0] != "en" {
		t.Fatalf("US langs = %v", us.Languages())
	}
	if len(us.TimeZones()) < 10 {
		t.Fatalf("US timezones = %d", len(us.TimeZones()))
	}

	undef := CountryByCode2("**")
	if undef.Currency() != nil || undef.Languages() != nil || undef.TimeZones() != nil {
		t.Fatalf("undefined lists must be nil")
	}
}

func TestContinentByCode2(t *testing.T) {
	eu := ContinentByCode2("EU")
	if eu == nil || eu.Name != "Europe" {
		t.Fatalf("EU: %+v", eu)
	}
	if ContinentByCode2("XX") != nil || ContinentByCode2("") != nil {
		t.Fatalf("unknown continent must be nil")
	}
}

func TestCountryMarshalJSON(t *testing.T) {
	data, err := CountryByCode2("AD").MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	s := string(data)
	for _, want := range []string{`"cc":"AD"`, `"ccc":"AND"`, `"currency":["EUR"]`, `"continent":"EU"`} {
		if !strings.Contains(s, want) {
			t.Fatalf("json missing %s: %s", want, s)
		}
	}
}

func TestUTCOffset(t *testing.T) {
	ad := CountryByCode2("AD")
	zones := ad.TimeZones()
	if len(zones) != 1 || zones[0].ZoneName != "Europe/Andorra" || zones[0].UTCOffset != 60 {
		t.Fatalf("Andorra zones = %+v", zones)
	}
	if ad.AvgUTCOffset != zones[0].UTCOffset {
		t.Fatalf("Andorra AvgUTCOffset = %d, want %d", ad.AvgUTCOffset, zones[0].UTCOffset)
	}

	us := CountryByCode2("US")
	var sum int
	var ny int16
	foundNY := false
	for _, z := range us.TimeZones() {
		sum += int(z.UTCOffset)
		if z.ZoneName == "America/New_York" {
			ny = z.UTCOffset
			foundNY = true
		}
	}
	if !foundNY || ny != -300 {
		t.Fatalf("America/New_York = %d, found=%v", ny, foundNY)
	}
	wantAvg := int16(math.Round(float64(sum) / float64(len(us.TimeZones()))))
	if us.AvgUTCOffset != wantAvg {
		t.Fatalf("US AvgUTCOffset = %d, want %d", us.AvgUTCOffset, wantAvg)
	}

	np := CountryByCode2("NP").TimeZones()
	if len(np) != 1 || np[0].ZoneName != "Asia/Katmandu" || np[0].UTCOffset != 345 {
		t.Fatalf("Kathmandu = %+v", np)
	}

	var adelaide int16
	foundAdelaide := false
	for _, z := range CountryByCode2("AU").TimeZones() {
		if z.ZoneName == "Australia/Adelaide" {
			adelaide = z.UTCOffset
			foundAdelaide = true
		}
	}
	if !foundAdelaide || adelaide != 570 {
		t.Fatalf("Australia/Adelaide = %d, found=%v", adelaide, foundAdelaide)
	}

	if got := RegionByCode("US-CA").UTCOffset; got != -480 {
		t.Fatalf("US-CA UTCOffset = %d, want -480", got)
	}
	if got := RegionByCode("AD-07").UTCOffset; got != 60 {
		t.Fatalf("AD-07 UTCOffset = %d, want 60", got)
	}
	if CountryByCode2("**").AvgUTCOffset != 0 || RegionByCode("??-XX").UTCOffset != 0 {
		t.Fatal("unknown country and region offsets must be 0")
	}
}

func TestCountriesIndex(t *testing.T) {
	if len(Countries) != 255 {
		t.Fatalf("len(Countries) = %d", len(Countries))
	}
	for i := range Countries {
		if int(Countries[i].ID) != i {
			t.Fatalf("id/index mismatch at %d", i)
		}
		if CountryByCode2(Countries[i].ISO2()) != &Countries[i] {
			t.Fatalf("ISO2 roundtrip failed for %s", Countries[i].ISO2())
		}
		if iso3 := Countries[i].ISO3(); iso3 != "" {
			if CountryByCode3(iso3) != &Countries[i] {
				t.Fatalf("ISO3 roundtrip failed for %s", iso3)
			}
		}
	}
}
