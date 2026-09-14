package gogeo

import (
	"runtime"
	"testing"
)

var (
	benchCountry *Country
	benchCode    Code2
	benchString  string
	benchBytes   []byte
)

func BenchmarkCountryByCode2(b *testing.B) {
	b.ReportAllocs()
	var c *Country
	for i := 0; i < b.N; i++ {
		c = CountryByCode2("US")
	}
	benchCountry = c
}

func BenchmarkCountryByCode2Bytes(b *testing.B) {
	code := Code2{'U', 'S'}
	b.ReportAllocs()
	var c *Country
	for i := 0; i < b.N; i++ {
		c = CountryByCode2Bytes(code)
	}
	benchCountry = c
}

func BenchmarkCountryByCode3(b *testing.B) {
	b.ReportAllocs()
	var c *Country
	for i := 0; i < b.N; i++ {
		c = CountryByCode3("USA")
	}
	benchCountry = c
}

func BenchmarkCountryCode2ByStringISO2(b *testing.B) {
	b.ReportAllocs()
	var cc Code2
	for i := 0; i < b.N; i++ {
		cc = CountryCode2ByString("US")
	}
	benchCode = cc
}

func BenchmarkCountryCode2ByStringISO3(b *testing.B) {
	b.ReportAllocs()
	var cc Code2
	for i := 0; i < b.N; i++ {
		cc = CountryCode2ByString("USA")
	}
	benchCode = cc
}

func BenchmarkCode2ISO2(b *testing.B) {
	cc := Code2{'U', 'S'}
	b.ReportAllocs()
	var s string
	for i := 0; i < b.N; i++ {
		s = cc.ISO2()
	}
	benchString = s
}

func BenchmarkCode2ISO3(b *testing.B) {
	cc := Code2{'U', 'S'}
	b.ReportAllocs()
	var s string
	for i := 0; i < b.N; i++ {
		s = cc.ISO3()
	}
	benchString = s
}

func BenchmarkCode2MarshalJSON(b *testing.B) {
	cc := Code2{'U', 'S'}
	b.ReportAllocs()
	var data []byte
	for i := 0; i < b.N; i++ {
		data, _ = cc.MarshalJSON()
	}
	benchBytes = data
}

func BenchmarkCode2Value(b *testing.B) {
	cc := Code2{'U', 'S'}
	b.ReportAllocs()
	var s string
	for i := 0; i < b.N; i++ {
		v, _ := cc.Value()
		s, _ = v.(string)
	}
	benchString = s
}

func TestCode2Allocs(t *testing.T) {
	cc := Code2{'U', 'S'}
	if n := testing.AllocsPerRun(1000, func() { _ = CountryByCode2("US") }); n != 0 {
		t.Errorf("CountryByCode2 allocs = %v, want 0", n)
	}
	if n := testing.AllocsPerRun(1000, func() { _ = cc.ISO2() }); n != 0 {
		t.Errorf("ISO2 allocs = %v, want 0", n)
	}
	unknown := Code2{'-', '-'}
	if n := testing.AllocsPerRun(1000, func() { _ = unknown.ISO2() }); n != 0 {
		t.Errorf("ISO2(unknown) allocs = %v, want 0", n)
	}
	if n := testing.AllocsPerRun(1000, func() { _ = cc.ISO3() }); n != 0 {
		t.Errorf("ISO3 allocs = %v, want 0", n)
	}
	if n := testing.AllocsPerRun(1000, func() { _, _ = cc.MarshalJSON() }); n > 1 {
		t.Errorf("MarshalJSON allocs = %v, want <= 1", n)
	}
	if n := testing.AllocsPerRun(1000, func() { _, _ = cc.Value() }); n != 0 {
		t.Errorf("Value allocs = %v, want 0", n)
	}
}

func TestImportHeap(t *testing.T) {
	runtime.GC()
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	t.Logf("heap_alloc=%d heap_objects=%d sys=%d countries=%d", ms.HeapAlloc, ms.HeapObjects, ms.Sys, len(Countries))
}
