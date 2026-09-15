//
// @project GeniusRabbit 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package gogeo

// Continent description
type Continent struct {
	ID    uint8
	Code2 string
	Name  string
}

// Continents list (IDs are 1-based; Country.continent stores this ID).
var Continents = []Continent{
	{ID: 1, Code2: "EU", Name: "Europe"},
	{ID: 2, Code2: "AS", Name: "Asia"},
	{ID: 3, Code2: "AF", Name: "Africa"},
	{ID: 4, Code2: "OC", Name: "Oceania"},
	{ID: 5, Code2: "SA", Name: "South America"},
	{ID: 6, Code2: "NA", Name: "North America"},
	{ID: 7, Code2: "AN", Name: "Antarctica"},
}

// ContinentByCode2 returns a continent by ISO-like 2-letter code, or nil.
func ContinentByCode2(code string) *Continent {
	if len(code) != 2 {
		return nil
	}
	id := continentID(code[0], code[1])
	if id == 0 {
		return nil
	}
	return &Continents[id-1]
}

func continentID(b0, b1 byte) uint8 {
	if b0 < 'A' || b0 > 'Z' || b1 < 'A' || b1 > 'Z' {
		return 0
	}
	return continentIndex[b0-'A'][b1-'A']
}

var continentIndex = [26][26]uint8{
	0:  {5: 3, 13: 7, 18: 2}, // AF, AN, AS
	4:  {20: 1},             // EU
	13: {0: 6},              // NA
	14: {2: 4},              // OC
	18: {0: 5},              // SA
}
