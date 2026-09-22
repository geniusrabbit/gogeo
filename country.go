//
// @project GeniusRabbit 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

package gogeo

import "encoding/json"

//go:generate go run ./internal/cmd/gencountries -src data/countries.json -out country_data.go

// undefined country codes
const (
	UndefinedCountryCodeISO2 = "**"
	UndefinedCountryCodeISO3 = "***"
)

// TimeZone object
type TimeZone struct {
	ZoneName string
	Lon      float32
}

// Coordinates object
type Coordinates struct {
	Lat, Lon float32
}

// Country is a compact country record. Variable-length lists are stored in
// package-level packed tables and exposed via methods.
type Country struct {
	ID          uint8
	Code2       CountryCode2
	Name        string
	Native      string
	Capital     string
	Coordinates Coordinates

	continent uint8
	currN     uint8
	langN     uint8
	phoneN    uint8
	zoneN     uint8
	code3     [3]byte
	currOff   uint16
	langOff   uint16
	phoneOff  uint16
	zoneOff   uint16
}

// ISO2 returns the interned ISO-3166 alpha-2 code.
func (c *Country) ISO2() string {
	return countryISO2[c.ID]
}

// ISO3 returns the interned ISO-3166 alpha-3 code (empty if the entry has none).
func (c *Country) ISO3() string {
	return countryISO3[c.ID]
}

// Code3 returns the interned ISO-3166 alpha-3 code.
func (c *Country) Code3() string {
	return countryISO3[c.ID]
}

// Continent returns the interned continent code (EU, AS, ...) or empty.
func (c *Country) Continent() string {
	if c.continent == 0 {
		return ""
	}
	return Continents[c.continent-1].Code2
}

// Currency returns the packed currency codes for the country.
func (c *Country) Currency() []string {
	if c.currN == 0 {
		return nil
	}
	return currencyTable[c.currOff : c.currOff+uint16(c.currN)]
}

// Languages returns the packed language codes for the country.
func (c *Country) Languages() []string {
	if c.langN == 0 {
		return nil
	}
	return languageTable[c.langOff : c.langOff+uint16(c.langN)]
}

// Phones returns the packed phone prefixes for the country.
func (c *Country) Phones() []string {
	if c.phoneN == 0 {
		return nil
	}
	return phoneTable[c.phoneOff : c.phoneOff+uint16(c.phoneN)]
}

// TimeZones returns the packed timezone list for the country.
func (c *Country) TimeZones() []TimeZone {
	if c.zoneN == 0 {
		return nil
	}
	return timezoneTable[c.zoneOff : c.zoneOff+uint16(c.zoneN)]
}

type countryJSON struct {
	ID          uint8       `json:"id"`
	Code2       string      `json:"cc,omitempty"`
	Code3       string      `json:"ccc,omitempty"`
	Name        string      `json:"name,omitempty"`
	Native      string      `json:"native,omitempty"`
	Continent   string      `json:"continent,omitempty"`
	Capital     string      `json:"capital,omitempty"`
	Currency    []string    `json:"currency,omitempty"`
	Languages   []string    `json:"langs,omitempty"`
	Phones      []string    `json:"phones,omitempty"`
	Coordinates Coordinates `json:"corrd,omitempty"`
	TimeZones   []TimeZone  `json:"zones,omitempty"`
}

// MarshalJSON encodes Country using the historical field names.
func (c Country) MarshalJSON() ([]byte, error) {
	return json.Marshal(countryJSON{
		ID:          c.ID,
		Code2:       c.ISO2(),
		Code3:       c.ISO3(),
		Name:        c.Name,
		Native:      c.Native,
		Continent:   c.Continent(),
		Capital:     c.Capital,
		Currency:    c.Currency(),
		Languages:   c.Languages(),
		Phones:      c.Phones(),
		Coordinates: c.Coordinates,
		TimeZones:   c.TimeZones(),
	})
}
