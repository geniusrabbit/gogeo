//
// @project GeniusRabbit 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

package gogeo

// undefined country codes
const (
	UndefinedCountryCodeISO2 = "**"
	UndefinedCountryCodeISO3 = "***"
)

// TimeZone object
type TimeZone struct {
	ZoneName string
	Lon      float64
}

// Coordinates object
type Coordinates struct {
	Lat, Lon float64
}

// Country info
type Country struct {
	ID          int         `json:"id"`
	Code2       string      `json:"cc,omitempty"`
	Code3       string      `json:"ccc,omitempty"`
	Name        string      `json:"name,omitempty"`
	ShortName   string      `json:"short,omitempty"`
	Native      string      `json:"native,omitempty"`
	Continent   string      `json:"continent,omitempty"`
	Capital     string      `json:"capital,omitempty"`
	Currency    []string    `json:"currency,omitempty"`
	Languages   []string    `json:"langs,omitempty"`
	Phones      []string    `json:"phones,omitempty"`
	Coordinates Coordinates `json:"corrd,omitempty"`
	TimeZones   []TimeZone  `json:"zones,omitempty"`
}
