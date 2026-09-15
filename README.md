# gogeo

[![Test](https://github.com/geniusrabbit/gogeo/actions/workflows/test.yml/badge.svg)](https://github.com/geniusrabbit/gogeo/actions/workflows/test.yml)
[![Lint](https://github.com/geniusrabbit/gogeo/actions/workflows/lint.yml/badge.svg)](https://github.com/geniusrabbit/gogeo/actions/workflows/lint.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/geniusrabbit/gogeo)](https://goreportcard.com/report/github.com/geniusrabbit/gogeo)
[![Go Reference](https://pkg.go.dev/badge/github.com/geniusrabbit/gogeo.svg)](https://pkg.go.dev/github.com/geniusrabbit/gogeo)

Static country, continent, and ISO 3166-2 region data with O(1) country lookups, packed region tables, and `Code2` / `RegionCode` types for JSON and SQL.

License: Apache 2.0

## Install

```bash
go get github.com/geniusrabbit/gogeo
```

Requires Go 1.21+.

## Usage

Unknown ISO codes resolve to the undefined country (`**` / `***`), never `nil`.

```go
package main

import (
    "fmt"

    "github.com/geniusrabbit/gogeo"
)

func main() {
    us := gogeo.CountryByCode2("US")
    fmt.Println(us.Name, us.ISO3(), us.Continent(), us.Currency())

    // ISO-2 or ISO-3 → Code2
    cc := gogeo.CountryCode2ByString("USA")
    fmt.Println(cc.ISO2(), cc.ISO3()) // US USA

    // Zero-alloc lookup when you already have Code2
    _ = gogeo.CountryByCode2Bytes(gogeo.Code2{'D', 'E'})

    eu := gogeo.ContinentByCode2("EU")
    fmt.Println(eu.Name)

    ca := gogeo.RegionByCode("US-CA")
    fmt.Println(ca.Name, ca.Type(), ca.Country().Name) // California State United States

    // Official ISO uses GB-*; gogeo country code is UK
    eng := gogeo.RegionByCode("UK-ENG")
    fmt.Println(eng.Code(), eng.Name) // GB-ENG England

    near := gogeo.RegionByLatLng(42.51, 1.52)
    fmt.Println(near.Code(), near.Name)

    _ = us.NearestRegion(34.05, -118.25)
}
```

`Code2` and `RegionCode` implement `json.Marshaler` / `Unmarshaler`, `driver.Valuer`, and `sql.Scanner`:

```go
type User struct {
    Country gogeo.Code2      `json:"country" db:"country_code"`
    Region  gogeo.RegionCode `json:"region" db:"region_code"`
}
```

Lists on `Country` are methods over packed tables: `Currency()`, `Languages()`, `Phones()`, `TimeZones()`, `Regions()`.

The dataset includes GeoIP-style extras (`A1`, `A2`, `O1`, `AP`, `EU`) and `UK` for the United Kingdom (not `GB`). ISO 3166-2 codes stay official (`GB-ENG`); `RegionByCode` accepts both `GB-*` and `UK-*`.

Subdivision data is merged from [Debian iso-codes](https://salsa.debian.org/iso-codes-team/iso-codes) (LGPL-2.1-or-later) and [oodavid/iso-3166-2](https://github.com/oodavid/iso-3166-2) centroids. See [`data/NOTICE`](data/NOTICE) and [`data/regions.json`](data/regions.json). `Region.Name` is English when available.

## Development

Country and region tables are generated from [`data/countries.json`](data/countries.json) and merged [`data/regions.json`](data/regions.json):

```bash
make build-data   # or: make generate
```

| Target | Description |
|---|---|
| `make test` | Unit tests |
| `make lint` | `gofmt` check + `go vet` |
| `make fmt` | Format sources |
| `make merge-data` | Merge debian + oodavid into `data/regions.json` |
| `make build-data` | Merge sources and regenerate Go tables |
| `make export-data` | Dump countries or regions (`KIND=countries\|regions FORMAT=json\|csv`) |
| `make bench` | Benchmarks |
| `make cover` | Coverage report |
| `make tidy` | `go mod tidy` |
| `make clean` | Remove coverage artifacts |

## Benchmarks

`make bench` on darwin/arm64 (Apple M2 Ultra):

```
BenchmarkCountryByCode2-24               1.967 ns/op    0 B/op    0 allocs/op
BenchmarkCountryByCode2Bytes-24          1.892 ns/op    0 B/op    0 allocs/op
BenchmarkCountryByCode3-24               2.175 ns/op    0 B/op    0 allocs/op
BenchmarkCountryCode2ByStringISO2-24     2.682 ns/op    0 B/op    0 allocs/op
BenchmarkCountryCode2ByStringISO3-24     2.937 ns/op    0 B/op    0 allocs/op
BenchmarkCode2ISO2-24                    2.315 ns/op    0 B/op    0 allocs/op
BenchmarkCode2ISO3-24                    2.103 ns/op    0 B/op    0 allocs/op
BenchmarkCode2MarshalJSON-24             4.318 ns/op    4 B/op    1 allocs/op
BenchmarkCode2Value-24                   2.348 ns/op    0 B/op    0 allocs/op
```
