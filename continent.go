//
// @project GeniusRabbit 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package gogeo

// Continent description
type Continent struct {
	ID    int
	Code2 string
	Name  string
}

// Continents list
var Continents = []Continent{
	{ID: 1, Code2: "EU", Name: "Europe"},
	{ID: 2, Code2: "AS", Name: "Asia"},
	{ID: 3, Code2: "AF", Name: "Africa"},
	{ID: 4, Code2: "OC", Name: "Oceania"},
	{ID: 5, Code2: "SA", Name: "South America"},
	{ID: 6, Code2: "NA", Name: "North America"},
	{ID: 7, Code2: "AN", Name: "Antarctica"},
}
