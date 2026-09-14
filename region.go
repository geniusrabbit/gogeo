package gogeo

//go:generate go run ./internal/cmd/genregions -src data/iso_3166-2.json -countries data/countries.json -out region_data.go

// Undefined ISO 3166-2 code.
const UndefinedRegionISO3166 = "**-**"

// Region is a compact ISO 3166-2 subdivision record.
type Region struct {
	ID      uint16
	Name    string
	parent  uint16
	country uint8
	typeID  uint8
	sufN    uint8
	suffix  [3]byte
}

// Code returns the interned ISO 3166-2 code (for example "US-CA").
func (r *Region) Code() string {
	return regionCodes[r.ID]
}

// Type returns the interned subdivision type (State, Province, Parish, ...).
func (r *Region) Type() string {
	return regionTypes[r.typeID]
}

// Country returns the parent country. GB subdivisions resolve to UK.
// Never nil.
func (r *Region) Country() *Country {
	return &countries[r.country]
}

// Parent returns the parent subdivision, or nil if the region is top-level.
func (r *Region) Parent() *Region {
	if r.parent == 0 {
		return nil
	}
	return &regions[r.parent]
}

// RegionCode returns the JSON/SQL code wrapper for this region.
func (r *Region) RegionCode() RegionCode {
	return RegionCode(r.ID)
}

// Regions returns ISO 3166-2 subdivisions for the country.
func (c *Country) Regions() []Region {
	rng := countryRegionRange[c.ID]
	if rng.n == 0 {
		return nil
	}
	return regions[rng.off : rng.off+rng.n]
}
