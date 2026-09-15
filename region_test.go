package gogeo

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestRegionByCode(t *testing.T) {
	ca := RegionByCode("US-CA")
	if ca == nil || ca.Code() != "US-CA" || ca.Name != "California" || ca.Type() != "State" {
		t.Fatalf("US-CA: code=%q name=%q type=%q", ca.Code(), ca.Name, ca.Type())
	}
	if ca.Country().ISO2() != "US" || ca.Country().Name != "United States" {
		t.Fatalf("US-CA country = %q %q", ca.Country().ISO2(), ca.Country().Name)
	}
	if ca.Parent() != nil {
		t.Fatalf("US-CA parent = %v", ca.Parent())
	}

	ad := RegionByCode("AD-07")
	if ad.Name != "Andorra la Vella" || ad.Type() != "Parish" || ad.Country().ISO2() != "AD" {
		t.Fatalf("AD-07: %+v type=%q country=%q", ad, ad.Type(), ad.Country().ISO2())
	}

	undef := RegionByCode("??-XX")
	if undef == nil || undef.Code() != UndefinedRegionISO3166 {
		t.Fatalf("unknown must be undefined, got %+v", undef)
	}
	if RegionByCode("") != undef || RegionByCode("US") != undef || RegionByCode("US-") != undef {
		t.Fatalf("invalid codes must be undefined")
	}
	if RegionByCode("us-ca") != undef {
		t.Fatalf("lookup must be case-sensitive")
	}
}

func TestRegionGBUKAlias(t *testing.T) {
	gb := RegionByCode("GB-ENG")
	uk := RegionByCode("UK-ENG")
	if gb != uk {
		t.Fatalf("GB-ENG and UK-ENG must be the same region")
	}
	if gb.Name != "England" || gb.Code() != "GB-ENG" {
		t.Fatalf("England: name=%q code=%q", gb.Name, gb.Code())
	}
	if gb.Country().ISO2() != "UK" {
		t.Fatalf("GB-ENG country = %q, want UK", gb.Country().ISO2())
	}
}

func TestRegionParent(t *testing.T) {
	r := RegionByCode("AZ-KAN")
	p := r.Parent()
	if p == nil || p.Code() != "AZ-NX" {
		t.Fatalf("AZ-KAN parent = %v", p)
	}
}

func TestCountryRegions(t *testing.T) {
	us := CountryByCode2("US").Regions()
	if len(us) < 50 {
		t.Fatalf("US regions = %d", len(us))
	}
	found := false
	for i := range us {
		if us[i].Code() == "US-CA" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("US.Regions() missing US-CA")
	}

	uk := CountryByCode2("UK").Regions()
	if len(uk) == 0 {
		t.Fatal("UK.Regions() empty")
	}
	if CountryByCode2("**").Regions() != nil {
		t.Fatal("undefined country must have no regions")
	}
}

func TestRegionByID(t *testing.T) {
	ca := RegionByCode("US-CA")
	if RegionByID(ca.ID) != ca {
		t.Fatal("RegionByID mismatch")
	}
	if RegionByID(9999).Code() != UndefinedRegionISO3166 {
		t.Fatal("out-of-range ID must be undefined")
	}
	if int(Regions[ca.ID].ID) != int(ca.ID) {
		t.Fatal("Regions index != ID")
	}
}

func TestRegionCodeJSONSQL(t *testing.T) {
	rc := RegionCodeByString("US-CA")
	if rc.ISO3166() != "US-CA" || rc.Country() != (Code2{'U', 'S'}) {
		t.Fatalf("RegionCode US-CA = %q country=%q", rc.ISO3166(), rc.Country())
	}
	if rc.Region() != RegionByCode("US-CA") {
		t.Fatal("RegionCode.Region mismatch")
	}

	data, err := json.Marshal(rc)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `"US-CA"` {
		t.Fatalf("marshal = %s", data)
	}

	var got RegionCode
	if err := json.Unmarshal([]byte(`"GB-ENG"`), &got); err != nil {
		t.Fatal(err)
	}
	if got.ISO3166() != "GB-ENG" || got.Country() != (Code2{'U', 'K'}) {
		t.Fatalf("unmarshal GB-ENG = %q country=%q", got.ISO3166(), got.Country())
	}

	if err := json.Unmarshal([]byte(`"nope"`), &got); err != nil {
		t.Fatal(err)
	}
	if got != UndefinedRegionCode {
		t.Fatalf("invalid json code = %q", got)
	}

	if err := got.Scan("AD-07"); err != nil {
		t.Fatal(err)
	}
	if got.ISO3166() != "AD-07" {
		t.Fatalf("scan = %q", got)
	}
	if err := got.Scan(123); !errors.Is(err, ErrRegionCodeInvalidScanType) {
		t.Fatalf("scan int err = %v", err)
	}
	v, err := RegionCodeByString("US-CA").Value()
	if err != nil || v != "US-CA" {
		t.Fatalf("value = %v %v", v, err)
	}
}

func TestRegionCoordinates(t *testing.T) {
	ad := RegionByCode("AD-07")
	if !ad.HasCoordinates() || ad.Coordinates.Lat == 0 || ad.Coordinates.Lon == 0 {
		t.Fatalf("AD-07 coords = %+v has=%v", ad.Coordinates, ad.HasCoordinates())
	}
	ca := RegionByCode("US-CA")
	if !ca.HasCoordinates() {
		t.Fatal("US-CA should have coordinates")
	}
	if names := ca.Names(); len(names) == 0 || names[0] != ca.Name {
		t.Fatalf("US-CA Names() = %v, Name=%q", names, ca.Name)
	}
	if ca.Name != "California" {
		t.Fatalf("US-CA Name = %q, want English default", ca.Name)
	}
}

func TestRegionByLatLng(t *testing.T) {
	ad := RegionByCode("AD-07")
	got := RegionByLatLng(ad.Coordinates.Lat, ad.Coordinates.Lon)
	if got.Country().ISO2() != "AD" {
		t.Fatalf("nearest Andorra = %s %s", got.Code(), got.Country().ISO2())
	}

	ca := RegionByCode("US-CA")
	near := CountryByCode2("US").NearestRegion(ca.Coordinates.Lat, ca.Coordinates.Lon)
	if near.Country().ISO2() != "US" {
		t.Fatalf("US nearest left the country: %s %s", near.Code(), near.Country().ISO2())
	}
	if !near.HasCoordinates() {
		t.Fatal("US nearest has no coordinates")
	}

	if RegionByCode("**").HasCoordinates() {
		t.Fatal("undefined region must not have coordinates")
	}
	if CountryByCode2("**").NearestRegion(0, 0).Code() != UndefinedRegionISO3166 {
		t.Fatal("undefined country nearest must be undefined")
	}
}

func TestRegionAllocs(t *testing.T) {
	if n := testing.AllocsPerRun(1000, func() { _ = RegionByCode("US-CA") }); n != 0 {
		t.Errorf("RegionByCode allocs = %v, want 0", n)
	}
	if n := testing.AllocsPerRun(1000, func() { _ = CountryByCode2("US").Regions() }); n != 0 {
		t.Errorf("Country.Regions allocs = %v, want 0", n)
	}
	rc := RegionCodeByString("US-CA")
	if n := testing.AllocsPerRun(1000, func() { _ = rc.ISO3166() }); n != 0 {
		t.Errorf("RegionCode.ISO3166 allocs = %v, want 0", n)
	}
	if n := testing.AllocsPerRun(1000, func() { _, _ = rc.MarshalJSON() }); n > 1 {
		t.Errorf("RegionCode.MarshalJSON allocs = %v, want <= 1", n)
	}
	ad := RegionByCode("AD-07")
	if n := testing.AllocsPerRun(200, func() { _ = RegionByLatLng(ad.Coordinates.Lat, ad.Coordinates.Lon) }); n != 0 {
		t.Errorf("RegionByLatLng allocs = %v, want 0", n)
	}
}

var benchRegion *Region

func BenchmarkRegionByCode(b *testing.B) {
	b.ReportAllocs()
	var r *Region
	for i := 0; i < b.N; i++ {
		r = RegionByCode("US-CA")
	}
	benchRegion = r
}

func BenchmarkCountryRegions(b *testing.B) {
	us := CountryByCode2("US")
	b.ReportAllocs()
	var list []Region
	for i := 0; i < b.N; i++ {
		list = us.Regions()
	}
	if len(list) > 0 {
		benchRegion = &list[0]
	}
}

func BenchmarkRegionByLatLng(b *testing.B) {
	ad := RegionByCode("AD-07")
	lat, lng := ad.Coordinates.Lat, ad.Coordinates.Lon
	b.ReportAllocs()
	var r *Region
	for i := 0; i < b.N; i++ {
		r = RegionByLatLng(lat, lng)
	}
	benchRegion = r
}

func BenchmarkRegionCodeMarshalJSON(b *testing.B) {
	rc := RegionCodeByString("US-CA")
	b.ReportAllocs()
	var data []byte
	for i := 0; i < b.N; i++ {
		data, _ = rc.MarshalJSON()
	}
	benchBytes = data
}
