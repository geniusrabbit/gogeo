//
// @project GeniusRabbit 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package gogeo

var (
	// countryMapCode2 code2: object
	countryMapCode2 = map[string]*Country{}

	// countryMapCode3 code3: object
	countryMapCode3 = map[string]*Country{}
)

func init() {
	for _, c := range Countries {
		cc := new(Country)
		*cc = c
		countryMapCode2[c.Code2] = cc
		countryMapCode3[c.Code3] = cc
	}
}

// CountryByCode2 object
func CountryByCode2(code string) *Country {
	c, _ := countryMapCode2[code]
	return c
}

// CountryByCode3 object
func CountryByCode3(code string) *Country {
	c, _ := countryMapCode3[code]
	return c
}
