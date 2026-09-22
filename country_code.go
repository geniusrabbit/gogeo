package gogeo

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// UndefinedCountryCode2 for undefined codename
var UndefinedCountryCode2 = CountryCode2{UndefinedCountryCodeISO2[0], UndefinedCountryCodeISO2[1]}

// CountryCode2ByString returns code by string code ISO-2 or ISO-3.
func CountryCode2ByString(code string) CountryCode2 {
	switch len(code) {
	case 2:
		return CountryByCode2(code).Code2
	case 3:
		return CountryByCode3(code).Code2
	default:
		return UndefinedCountryCode2
	}
}

// Error list...
var (
	ErrCode2InvalidScanType  = errors.New("[gogeo.code2] invalid scan type, supports only bytes and string")
	ErrCode2InvalidValueSize = errors.New("[gogeo.code2] invalid value size, supports only two 2 chars")
)

// CountryCode2 implements a 2-letter country code.
type CountryCode2 [2]byte

func (cc CountryCode2) String() string {
	return cc.ISO2()
}

// ISO2 returns the interned two-letter code for a known country, or "**".
// Aliases resolve to the canonical code (UK → GB).
func (cc CountryCode2) ISO2() string {
	c := lookupCode2(cc[0], cc[1])
	if c.ID == 0 {
		return UndefinedCountryCodeISO2
	}
	return countryISO2[c.ID]
}

// ISO3 returns the interned three-letter code for the matching country.
func (cc CountryCode2) ISO3() string {
	return countryISO3[lookupCode2(cc[0], cc[1]).ID]
}

// Value implementation of sql driver.Valuer interface
func (cc CountryCode2) Value() (driver.Value, error) {
	return cc.ISO2(), nil
}

// Scan implementation of sql.Scanner interface
func (cc *CountryCode2) Scan(data interface{}) error {
	switch v := data.(type) {
	case []byte:
		if len(v) != 2 {
			return ErrCode2InvalidValueSize
		}
		*cc = CountryCode2{v[0], v[1]}
	case string:
		if len(v) != 2 {
			return ErrCode2InvalidValueSize
		}
		*cc = CountryCode2{v[0], v[1]}
	default:
		return ErrCode2InvalidScanType
	}
	return nil
}

// MarshalJSON implements json.Marshaler interface
func (cc CountryCode2) MarshalJSON() ([]byte, error) {
	return []byte{'"', cc[0], cc[1], '"'}, nil
}

// UnmarshalJSON implements json.Unmarshaler interface
func (cc *CountryCode2) UnmarshalJSON(data []byte) error {
	if n := len(data); n >= 2 && data[0] == '"' && data[n-1] == '"' {
		s := data[1 : n-1]
		if len(s) != 2 {
			*cc = UndefinedCountryCode2
			return nil
		}
		*cc = CountryCode2{s[0], s[1]}
		return nil
	}
	var code string
	if err := json.Unmarshal(data, &code); err != nil {
		return err
	}
	if code == UndefinedCountryCodeISO2 || len(code) != 2 {
		*cc = UndefinedCountryCode2
	} else {
		*cc = CountryCode2{code[0], code[1]}
	}
	return nil
}

var (
	_ driver.Valuer    = UndefinedCountryCode2
	_ sql.Scanner      = &UndefinedCountryCode2
	_ json.Marshaler   = UndefinedCountryCode2
	_ json.Unmarshaler = &UndefinedCountryCode2
)
