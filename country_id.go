package gogeo

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// UndefinedCountryID is the sentinel country index.
var UndefinedCountryID CountryID

// Error list for CountryID SQL scanning.
var (
	ErrCountryIDInvalidScanType  = errors.New("[gogeo.countryid] invalid scan type, supports only bytes and string")
	ErrCountryIDInvalidValueSize = errors.New("[gogeo.countryid] invalid value size, supports only two 2 chars")
)

// CountryID is a country index for JSON and SQL fields.
// On the wire it is the ISO-2 code, same as CountryCode2.
type CountryID uint8

// CountryIDByString parses an ISO-2 or ISO-3 code. Unknown values become undefined.
func CountryIDByString(code string) CountryID {
	switch len(code) {
	case 2:
		return CountryID(CountryByCode2(code).ID)
	case 3:
		return CountryID(CountryByCode3(code).ID)
	default:
		return UndefinedCountryID
	}
}

func (id CountryID) String() string {
	return id.ISO2()
}

// Object returns the country for this index, or the undefined country.
// Never nil.
func (id CountryID) Object() *Country {
	return CountryByID(uint8(id))
}

// Code returns the canonical ISO-2 code. UK resolves to GB.
func (id CountryID) Code() CountryCode2 {
	return id.Object().Code2
}

// ID returns the interned country index, or 0 for undefined.
func (id CountryID) ID() uint8 {
	return id.Object().ID
}

// ISO2 returns the interned two-letter code, or "**".
func (id CountryID) ISO2() string {
	return id.Object().ISO2()
}

// ISO3 returns the interned three-letter code for the matching country.
func (id CountryID) ISO3() string {
	return id.Object().ISO3()
}

// Value implementation of sql driver.Valuer.
func (id CountryID) Value() (driver.Value, error) {
	return id.ISO2(), nil
}

// Scan implementation of sql.Scanner.
func (id *CountryID) Scan(data interface{}) error {
	switch v := data.(type) {
	case []byte:
		if len(v) != 2 {
			return ErrCountryIDInvalidValueSize
		}
		*id = CountryID(CountryByCode2(string(v)).ID)
	case string:
		if len(v) != 2 {
			return ErrCountryIDInvalidValueSize
		}
		*id = CountryID(CountryByCode2(v).ID)
	default:
		return ErrCountryIDInvalidScanType
	}
	return nil
}

// MarshalJSON implements json.Marshaler.
func (id CountryID) MarshalJSON() ([]byte, error) {
	s := id.ISO2()
	return []byte{'"', s[0], s[1], '"'}, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (id *CountryID) UnmarshalJSON(data []byte) error {
	if n := len(data); n >= 2 && data[0] == '"' && data[n-1] == '"' {
		s := data[1 : n-1]
		if len(s) != 2 {
			*id = UndefinedCountryID
			return nil
		}
		*id = CountryID(CountryByCode2(string(s)).ID)
		return nil
	}
	var code string
	if err := json.Unmarshal(data, &code); err != nil {
		return err
	}
	if code == UndefinedCountryCodeISO2 || len(code) != 2 {
		*id = UndefinedCountryID
	} else {
		*id = CountryID(CountryByCode2(code).ID)
	}
	return nil
}

var (
	_ driver.Valuer    = UndefinedCountryID
	_ sql.Scanner      = &UndefinedCountryID
	_ json.Marshaler   = UndefinedCountryID
	_ json.Unmarshaler = &UndefinedCountryID
)
