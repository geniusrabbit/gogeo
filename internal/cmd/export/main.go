// Command export dumps gogeo countries or ISO 3166-2 regions as JSON or CSV.
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

func main() {
	kind := flag.String("kind", "countries", "countries | regions")
	format := flag.String("format", "json", "json | csv")
	flag.Parse()

	switch *kind {
	case "countries":
		if err := exportCountries(*format); err != nil {
			log.Fatal(err)
		}
	case "regions":
		if err := exportRegions(*format); err != nil {
			log.Fatal(err)
		}
	default:
		flag.Usage()
		os.Exit(2)
	}
}

func exportCountries(format string) error {
	switch format {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(gogeo.Countries)
	case "csv":
		wr := csv.NewWriter(os.Stdout)
		if err := wr.Write([]string{
			"id", "code2", "code3", "name", "native",
			"continent", "capital", "currency", "langs",
			"phones", "lat", "lon", "zones",
		}); err != nil {
			return err
		}
		for i := range gogeo.Countries {
			c := &gogeo.Countries[i]
			zones, err := json.Marshal(c.TimeZones())
			if err != nil {
				return err
			}
			if err := wr.Write([]string{
				strconv.Itoa(int(c.ID)), c.ISO2(), c.ISO3(), c.Name, c.Native,
				c.Continent(), c.Capital,
				strings.Join(c.Currency(), ","),
				strings.Join(c.Languages(), ","),
				strings.Join(c.Phones(), ","),
				fmt.Sprintf("%f", c.Coordinates.Lat),
				fmt.Sprintf("%f", c.Coordinates.Lon),
				string(zones),
			}); err != nil {
				return err
			}
		}
		wr.Flush()
		return wr.Error()
	default:
		return fmt.Errorf("unsupported format %q", format)
	}
}

func exportRegions(format string) error {
	switch format {
	case "json":
		out := make([]regionJSON, 0, len(gogeo.Regions))
		for i := range gogeo.Regions {
			out = append(out, newRegionJSON(&gogeo.Regions[i]))
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	case "csv":
		wr := csv.NewWriter(os.Stdout)
		if err := wr.Write([]string{
			"id", "code", "name", "names", "type", "country", "parent", "lat", "lng",
		}); err != nil {
			return err
		}
		for i := range gogeo.Regions {
			r := &gogeo.Regions[i]
			parent := ""
			if p := r.Parent(); p != nil {
				parent = p.Code()
			}
			lat, lng := "", ""
			if r.HasCoordinates() {
				lat = fmt.Sprintf("%f", r.Coordinates.Lat)
				lng = fmt.Sprintf("%f", r.Coordinates.Lon)
			}
			if err := wr.Write([]string{
				strconv.Itoa(int(r.ID)), r.Code(), r.Name, strings.Join(r.Names(), "|"),
				r.Type(), r.Country().ISO2(), parent, lat, lng,
			}); err != nil {
				return err
			}
		}
		wr.Flush()
		return wr.Error()
	default:
		return fmt.Errorf("unsupported format %q", format)
	}
}

type regionJSON struct {
	ID      uint16   `json:"id"`
	Code    string   `json:"code,omitempty"`
	Name    string   `json:"name,omitempty"`
	Names   []string `json:"names,omitempty"`
	Type    string   `json:"type,omitempty"`
	Country string   `json:"country,omitempty"`
	Parent  string   `json:"parent,omitempty"`
	Lat     *float32 `json:"lat,omitempty"`
	Lng     *float32 `json:"lng,omitempty"`
}

func newRegionJSON(r *gogeo.Region) regionJSON {
	out := regionJSON{
		ID:      r.ID,
		Code:    r.Code(),
		Name:    r.Name,
		Names:   r.Names(),
		Type:    r.Type(),
		Country: r.Country().ISO2(),
	}
	if p := r.Parent(); p != nil {
		out.Parent = p.Code()
	}
	if r.HasCoordinates() {
		lat, lng := r.Coordinates.Lat, r.Coordinates.Lon
		out.Lat, out.Lng = &lat, &lng
	}
	return out
}
