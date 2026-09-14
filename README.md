# gogeo

[![Test](https://github.com/geniusrabbit/gogeo/actions/workflows/test.yml/badge.svg)](https://github.com/geniusrabbit/gogeo/actions/workflows/test.yml)
[![Lint](https://github.com/geniusrabbit/gogeo/actions/workflows/lint.yml/badge.svg)](https://github.com/geniusrabbit/gogeo/actions/workflows/lint.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/geniusrabbit/gogeo)](https://goreportcard.com/report/github.com/geniusrabbit/gogeo)
[![Go Reference](https://pkg.go.dev/badge/github.com/geniusrabbit/gogeo.svg)](https://pkg.go.dev/github.com/geniusrabbit/gogeo)

Static country and continent reference data with O(1) ISO lookups and a `Code2` type for JSON and SQL.

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
}
```

`Code2` implements `json.Marshaler` / `Unmarshaler`, `driver.Valuer`, and `sql.Scanner`:

```go
type User struct {
    Country gogeo.Code2 `json:"country" db:"country_code"`
}
```

Lists on `Country` are methods over packed tables: `Currency()`, `Languages()`, `Phones()`, `TimeZones()`.

The dataset includes GeoIP-style extras (`A1`, `A2`, `O1`, `AP`, `EU`) and `UK` for the United Kingdom (not `GB`).

## Development

Country tables are generated from [`data/countries.json`](data/countries.json):

```bash
make build-data   # or: make generate
```

| Target | Description |
|---|---|
| `make test` | Unit tests |
| `make lint` | `gofmt` check + `go vet` |
| `make fmt` | Format sources |
| `make build-data` | Regenerate `country_data.go` |
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
