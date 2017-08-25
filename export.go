//
// @project GeniusRabbit 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

// +build ignore

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
			"id", "code2", "code3", "name", "short", "native",
			"continent", "campital", "currency", "langs",
			"phones", "lat", "lon", "zones",
		})
		for _, c := range gogeo.Countries {
			data, _ := json.Marshal(c.TimeZones)
			wr.Write([]string{
				strconv.Itoa(c.ID), c.Code2, c.Code3, c.Name, c.ShortName, c.Native, c.Continent, c.Capital,
				strings.Join(c.Currency, ","), strings.Join(c.Languages, ","), strings.Join(c.Phones, ","),
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
