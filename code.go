package gogeo

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// UndefinedCountryCode2 for undefined codename
var UndefinedCountryCode2 = Code2{UndefinedCountryCodeISO2[0], UndefinedCountryCodeISO2[1]}

// CountryCode2ByString returns code by string code ISO-2 or ISO-3
func CountryCode2ByString(code string) Code2 {
	var country *Country
	switch len(code) {
	case 2:
		country = CountryByCode2(code)
	case 3:
		country = CountryByCode3(code)
	}
	if country == nil {
		return UndefinedCountryCode2
	}
	return Code2{country.Code2[0], country.Code2[1]}
}

// Error list...
var (
	ErrCode2InvalidScanType  = errors.New("[gogeo.code2] invalid scan type, supports only bytes and string")
	ErrCode2InvalidValueSize = errors.New("[gogeo.code2] invalid value size, supports only two 2 chars")
)

// Code2 implements basic
type Code2 [2]byte

func (cc Code2) String() string {
	return cc.ISO2()
}

// ISO2 returns the string code with two simbols
func (cc Code2) ISO2() string {
	return string(cc[:])
}

// ISO3 returns the string code with three simbols
func (cc Code2) ISO3() string {
	if country := CountryByCode2(cc.ISO2()); country != nil {
		return country.Code3
	}
	return UndefinedCountryCodeISO3
}

// Value implementation of sql driver.Valuer interface
func (cc Code2) Value() (driver.Value, error) {
	return []byte(cc[:]), nil
}

// Scan implementation of sql.Scanner interface
func (cc *Code2) Scan(data interface{}) error {
	switch v := data.(type) {
	case []byte:
		if len(v) != 2 {
			*cc = Code2{v[0], v[1]}
		} else {
			return ErrCode2InvalidValueSize
		}
	case string:
		if len(v) != 2 {
			*cc = Code2{v[0], v[1]}
		} else {
			return ErrCode2InvalidValueSize
		}
	default:
		return ErrCode2InvalidScanType
	}
	return nil
}

// MarshalJSON implements json.Marshaler interface
func (cc Code2) MarshalJSON() ([]byte, error) {
	return json.Marshal(cc.ISO2())
}

// UnmarshalJSON implements json.Unmarshaler interface
func (cc *Code2) UnmarshalJSON(data []byte) error {
	var code string
	if err := json.Unmarshal(data, &code); err != nil {
		return err
	}
	if code == UndefinedCountryCodeISO2 || len(code) != 2 {
		*cc = UndefinedCountryCode2
	} else {
		*cc = Code2{code[0], code[1]}
	}
	return nil
}

var (
	_ driver.Valuer    = &UndefinedCountryCode2
	_ sql.Scanner      = &UndefinedCountryCode2
	_ json.Marshaler   = &UndefinedCountryCode2
	_ json.Unmarshaler = &UndefinedCountryCode2
)
