//
// @project GeniusRabbit 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

package gogeo

// Countries is the public list of all countries. Index equals Country.ID.
var Countries = countries[:]

// CountryByID returns the country with the given ID, or the undefined country.
func CountryByID(id uint8) *Country {
	if int(id) >= len(countries) {
		return &countries[0]
	}
	return &countries[id]
}

// CountryByCode2 returns the country for an ISO-2 code.
// Unknown codes resolve to the undefined country (never nil).
func CountryByCode2(code string) *Country {
	if len(code) != 2 {
		return &countries[0]
	}
	return lookupCode2(code[0], code[1])
}

// CountryByCode2Bytes is the allocation-free ISO-2 lookup.
// Unknown codes resolve to the undefined country (never nil).
func CountryByCode2Bytes(code Code2) *Country {
	return lookupCode2(code[0], code[1])
}

// CountryByCode3 returns the country for an ISO-3 code.
// Empty and unknown codes resolve to the undefined country (never nil).
func CountryByCode3(code string) *Country {
	if len(code) != 3 {
		return &countries[0]
	}
	return lookupCode3(code[0], code[1], code[2])
}

func lookupCode2(b0, b1 byte) *Country {
	if b0 >= 'A' && b0 <= 'Z' && b1 >= 'A' && b1 <= 'Z' {
		return &countries[code2Index[b0-'A'][b1-'A']]
	}
	return lookupCode2Special(b0, b1)
}

func lookupCode3(b0, b1, b2 byte) *Country {
	if b0 >= 'A' && b0 <= 'Z' && b1 >= 'A' && b1 <= 'Z' && b2 >= 'A' && b2 <= 'Z' {
		return &countries[code3Index[b0-'A'][b1-'A'][b2-'A']]
	}
	return lookupCode3Special(b0, b1, b2)
}
