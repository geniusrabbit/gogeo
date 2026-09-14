//
// @project GeniusRabbit 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

//go:build ignore

package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/geniusrabbit/gogeo"
)

var (
	flagFormat = flag.String("format", "json", "Type of export format json | csv")
)

func main() {
	flag.Parse()

	switch *flagFormat {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(gogeo.Countries)
	case "csv":
		wr := csv.NewWriter(os.Stdout)
		wr.Write([]string{
			"id", "code2", "code3", "name", "native",
			"continent", "capital", "currency", "langs",
			"phones", "lat", "lon", "zones",
		})
		for i := range gogeo.Countries {
			c := &gogeo.Countries[i]
			data, _ := json.Marshal(c.TimeZones())
			wr.Write([]string{
				strconv.Itoa(int(c.ID)), c.ISO2(), c.ISO3(), c.Name, c.Native, c.Continent(), c.Capital,
				strings.Join(c.Currency(), ","), strings.Join(c.Languages(), ","), strings.Join(c.Phones(), ","),
				fmt.Sprintf("%f", c.Coordinates.Lat), fmt.Sprintf("%f", c.Coordinates.Lon), string(data),
			})
		}
		wr.Flush()

		if err := wr.Error(); err != nil {
			log.Fatal(err)
		}
	default:
		flag.Usage()
	}
}
