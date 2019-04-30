//
// @project GeniusRabbit 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

package gogeo

var (
	// countryMapCode2 code2: object
	countryMapCode2 = map[string]*Country{}

	// countryMapCode3 code3: object
	countryMapCode3 = map[string]*Country{}
)

func init() {
	for i := range Countries {
		countryMapCode2[Countries[i].Code2] = &Countries[i]
		countryMapCode3[Countries[i].Code3] = &Countries[i]
	}
}

// CountryByCode2 object
func CountryByCode2(code string) *Country {
	if c, _ := countryMapCode2[code]; c != nil {
		return c
	}
	c, _ := countryMapCode2[UndefinedCountryCodeISO2]
	return c
}

// CountryByCode3 object
func CountryByCode3(code string) *Country {
	if c, _ := countryMapCode3[code]; c != nil {
		return c
	}
	c, _ := countryMapCode2[UndefinedCountryCodeISO2]
	return c
}
