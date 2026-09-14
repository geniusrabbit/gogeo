package gogeo

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// UndefinedRegionCode is the sentinel ISO 3166-2 code.
var UndefinedRegionCode = RegionCode(0)

// Error list for RegionCode SQL scanning.
var (
	ErrRegionCodeInvalidScanType = errors.New("[gogeo.region] invalid scan type, supports only bytes and string")
)

// RegionCode is an ISO 3166-2 subdivision code for JSON and SQL fields.
type RegionCode uint16

func (rc RegionCode) String() string {
	return rc.ISO3166()
}

// ISO3166 returns the interned ISO 3166-2 code, or "**-**".
func (rc RegionCode) ISO3166() string {
	return regionCodes[rc.clamp()]
}

// Country returns the library country code (GB subdivisions → UK).
func (rc RegionCode) Country() Code2 {
	return regions[rc.clamp()].Country().Code2
}

// Region returns the matching region (never nil).
func (rc RegionCode) Region() *Region {
	return RegionByID(uint16(rc.clamp()))
}

// RegionCodeByString parses an ISO 3166-2 code. Unknown values become undefined.
func RegionCodeByString(code string) RegionCode {
	return RegionCode(RegionByCode(code).ID)
}

func (rc RegionCode) clamp() RegionCode {
	if int(rc) >= len(regions) {
		return 0
	}
	return rc
}

// Value implementation of sql driver.Valuer.
func (rc RegionCode) Value() (driver.Value, error) {
	return rc.ISO3166(), nil
}

// Scan implementation of sql.Scanner.
func (rc *RegionCode) Scan(data interface{}) error {
	switch v := data.(type) {
	case []byte:
		*rc = RegionCodeByString(string(v))
	case string:
		*rc = RegionCodeByString(v)
	default:
		return ErrRegionCodeInvalidScanType
	}
	return nil
}

// MarshalJSON implements json.Marshaler.
func (rc RegionCode) MarshalJSON() ([]byte, error) {
	s := rc.ISO3166()
	out := make([]byte, 0, len(s)+2)
	out = append(out, '"')
	out = append(out, s...)
	return append(out, '"'), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (rc *RegionCode) UnmarshalJSON(data []byte) error {
	if n := len(data); n >= 2 && data[0] == '"' && data[n-1] == '"' {
		*rc = RegionCodeByString(string(data[1 : n-1]))
		return nil
	}
	var code string
	if err := json.Unmarshal(data, &code); err != nil {
		return err
	}
	*rc = RegionCodeByString(code)
	return nil
}

var (
	_ driver.Valuer    = UndefinedRegionCode
	_ sql.Scanner      = &UndefinedRegionCode
	_ json.Marshaler   = UndefinedRegionCode
	_ json.Unmarshaler = &UndefinedRegionCode
)
