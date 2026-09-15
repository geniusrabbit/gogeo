// Command genregions builds packed ISO 3166-2 tables from data/regions.json.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"os"
	"sort"
	"strconv"
	"strings"
)

type regionsFile struct {
	Regions []regionJSON `json:"regions"`
}

type regionJSON struct {
	Code   string            `json:"code"`
	Name   string            `json:"name"`
	Names  map[string]string `json:"names"`
	Type   string            `json:"type"`
	Parent string            `json:"parent"`
	Lat    *float64          `json:"lat"`
	Lng    *float64          `json:"lng"`
}

type countryJSON struct {
	ID    uint8  `json:"id"`
	Code2 string `json:"cc"`
}

func main() {
	src := flag.String("src", "data/regions.json", "merged regions JSON")
	countries := flag.String("countries", "data/countries.json", "gogeo countries JSON")
	out := flag.String("out", "region_data.go", "output Go file")
	flag.Parse()

	items, err := loadRegions(*src)
	if err != nil {
		fatal(err)
	}
	ccToID, err := loadCountries(*countries)
	if err != nil {
		fatal(err)
	}
	generated, err := build(*src, items, ccToID)
	if err != nil {
		fatal(err)
	}
	formatted, err := format.Source(generated)
	if err != nil {
		_, _ = os.Stderr.Write(generated)
		fatal(err)
	}
	if err := os.WriteFile(*out, formatted, 0644); err != nil {
		fatal(err)
	}
}

func loadRegions(path string) ([]regionJSON, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var file regionsFile
	if err := json.Unmarshal(raw, &file); err != nil {
		return nil, err
	}
	if len(file.Regions) == 0 {
		return nil, fmt.Errorf("no regions in %s", path)
	}
	return file.Regions, nil
}

func loadCountries(path string) (map[string]uint8, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var list []countryJSON
	if err := json.Unmarshal(raw, &list); err != nil {
		return nil, err
	}
	out := make(map[string]uint8, len(list))
	for _, c := range list {
		if c.Code2 != "" {
			out[c.Code2] = c.ID
		}
	}
	if uk, ok := out["UK"]; ok {
		out["GB"] = uk
	}
	return out, nil
}

type packedRegion struct {
	id       uint16
	name     string
	alts     []string
	code     string
	cc       string
	suffix   string
	typ      string
	parent   string
	country  uint8
	typeID   uint8
	parentID uint16
	lat, lng float32
	hasCoord bool
	nameOff  uint16
	nameN    uint8
}

func build(srcPath string, items []regionJSON, ccToID map[string]uint8) ([]byte, error) {
	types := map[string]uint8{"": 0}
	typeList := []string{""}

	byCC := map[string][]packedRegion{}
	seenCode := map[string]bool{}
	for _, item := range items {
		cc, suf, err := splitCode(item.Code)
		if err != nil {
			return nil, err
		}
		if seenCode[item.Code] {
			return nil, fmt.Errorf("duplicate region code %q", item.Code)
		}
		seenCode[item.Code] = true
		if _, ok := types[item.Type]; !ok {
			if len(types) > 255 {
				return nil, fmt.Errorf("too many region types")
			}
			types[item.Type] = uint8(len(typeList))
			typeList = append(typeList, item.Type)
		}
		name := item.Name
		if name == "" {
			if en := item.Names["en"]; en != "" {
				name = en
			} else {
				for _, v := range item.Names {
					if v != "" {
						name = v
						break
					}
				}
			}
		}
		alts := uniqueAlts(name, item.Names)
		pr := packedRegion{
			name:    name,
			alts:    alts,
			code:    item.Code,
			cc:      cc,
			suffix:  suf,
			typ:     item.Type,
			parent:  item.Parent,
			country: ccToID[cc],
			typeID:  types[item.Type],
		}
		if item.Lat != nil && item.Lng != nil && (*item.Lat != 0 || *item.Lng != 0) {
			pr.lat = float32(*item.Lat)
			pr.lng = float32(*item.Lng)
			pr.hasCoord = true
		}
		byCC[cc] = append(byCC[cc], pr)
	}

	ccs := make([]string, 0, len(byCC))
	for cc := range byCC {
		ccs = append(ccs, cc)
	}
	sort.Strings(ccs)

	packed := make([]packedRegion, 0, len(items)+1)
	packed = append(packed, packedRegion{
		id:     0,
		name:   "Undefined",
		code:   "**-**",
		cc:     "**",
		suffix: "**",
	})
	for _, cc := range ccs {
		list := byCC[cc]
		sort.Slice(list, func(i, j int) bool { return list[i].suffix < list[j].suffix })
		for i := 1; i < len(list); i++ {
			if list[i].suffix == list[i-1].suffix {
				return nil, fmt.Errorf("duplicate suffix %s-%s", cc, list[i].suffix)
			}
		}
		for i := range list {
			list[i].id = uint16(len(packed))
			packed = append(packed, list[i])
		}
		byCC[cc] = list
	}

	codeID := make(map[string]uint16, len(packed))
	for _, r := range packed {
		codeID[r.code] = r.id
	}
	for i := range packed {
		if packed[i].parent == "" {
			continue
		}
		pid, ok := codeID[packed[i].parent]
		if !ok {
			return nil, fmt.Errorf("region %s: unknown parent %q", packed[i].code, packed[i].parent)
		}
		packed[i].parentID = pid
	}

	var nameTable []string
	for i := range packed {
		off := uint16(len(nameTable))
		nameTable = append(nameTable, packed[i].name)
		nameTable = append(nameTable, packed[i].alts...)
		n := 1 + len(packed[i].alts)
		if n > 255 {
			return nil, fmt.Errorf("too many names for %s", packed[i].code)
		}
		packed[i].nameOff = off
		packed[i].nameN = uint8(n)
	}

	var geoIDs []uint16
	for _, r := range packed {
		if r.hasCoord {
			geoIDs = append(geoIDs, r.id)
		}
	}

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "// Code generated by internal/cmd/genregions from %s. DO NOT EDIT.\n", srcPath)
	buf.WriteString("package gogeo\n\n")

	writeStringTable(&buf, "regionCodes", codesOf(packed))
	writeStringTable(&buf, "regionTypes", typeList)
	writeStringTable(&buf, "regionNamesTable", nameTable)

	buf.WriteString("var regions = [...]Region{\n")
	for _, r := range packed {
		var suf [3]byte
		copy(suf[:], r.suffix)
		has := 0
		if r.hasCoord {
			has = 1
		}
		fmt.Fprintf(&buf, "\t{ID: %d, Name: %s, Coordinates: Coordinates{Lat: %s, Lon: %s}, parent: %d, country: %d, typeID: %d, sufN: %d, hasCoord: %d, suffix: [3]byte{%s}, nameOff: %d, nameN: %d},\n",
			r.id,
			strconv.Quote(r.name),
			fmtFloat(r.lat),
			fmtFloat(r.lng),
			r.parentID,
			r.country,
			r.typeID,
			len(r.suffix),
			has,
			bytesLit(suf[:]),
			r.nameOff,
			r.nameN,
		)
	}
	buf.WriteString("}\n\n")

	buf.WriteString("var geoRegionIDs = [...]uint16{\n")
	for i, id := range geoIDs {
		if i%16 == 0 {
			buf.WriteByte('\t')
		}
		fmt.Fprintf(&buf, "%d,", id)
		if i%16 == 15 || i == len(geoIDs)-1 {
			buf.WriteByte('\n')
		}
	}
	buf.WriteString("}\n\n")

	writeRanges(&buf, "iso2RegionRange", isoRanges(byCC))
	writeCountryRanges(&buf, ccToID, byCC)
	return buf.Bytes(), nil
}

func uniqueAlts(name string, names map[string]string) []string {
	seen := map[string]bool{name: true, "": true}
	var out []string
	keys := make([]string, 0, len(names))
	for k := range names {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := names[k]
		if seen[v] {
			continue
		}
		seen[v] = true
		out = append(out, v)
	}
	return out
}

func codesOf(packed []packedRegion) []string {
	out := make([]string, len(packed))
	for i, r := range packed {
		out[i] = r.code
	}
	return out
}

func splitCode(code string) (cc, suf string, err error) {
	i := strings.IndexByte(code, '-')
	if i != 2 || len(code) < 4 || len(code) > 6 {
		return "", "", fmt.Errorf("invalid ISO 3166-2 code %q", code)
	}
	cc, suf = code[:2], code[3:]
	if len(suf) == 0 || len(suf) > 3 {
		return "", "", fmt.Errorf("invalid ISO 3166-2 suffix in %q", code)
	}
	return cc, suf, nil
}

func isoRanges(byCC map[string][]packedRegion) [26][26]regionRangeLit {
	var out [26][26]regionRangeLit
	for cc, list := range byCC {
		if len(cc) != 2 || cc[0] < 'A' || cc[0] > 'Z' || cc[1] < 'A' || cc[1] > 'Z' {
			continue
		}
		out[cc[0]-'A'][cc[1]-'A'] = regionRangeLit{off: list[0].id, n: uint16(len(list))}
	}
	if gb := out['G'-'A']['B'-'A']; gb.n != 0 {
		out['U'-'A']['K'-'A'] = gb
	}
	return out
}

type regionRangeLit struct {
	off uint16
	n   uint16
}

func writeStringTable(buf *bytes.Buffer, name string, values []string) {
	fmt.Fprintf(buf, "var %s = [...]string{\n", name)
	for _, v := range values {
		fmt.Fprintf(buf, "\t%s,\n", strconv.Quote(v))
	}
	buf.WriteString("}\n\n")
}

func writeRanges(buf *bytes.Buffer, name string, az [26][26]regionRangeLit) {
	fmt.Fprintf(buf, "var %s = [26][26]regionRange{\n", name)
	for i := 0; i < 26; i++ {
		var parts []string
		for j := 0; j < 26; j++ {
			if az[i][j].n == 0 {
				continue
			}
			parts = append(parts, fmt.Sprintf("%d: {off: %d, n: %d}", j, az[i][j].off, az[i][j].n))
		}
		if len(parts) == 0 {
			continue
		}
		fmt.Fprintf(buf, "\t%d: {%s},\n", i, strings.Join(parts, ", "))
	}
	buf.WriteString("}\n\n")
}

func writeCountryRanges(buf *bytes.Buffer, ccToID map[string]uint8, byCC map[string][]packedRegion) {
	buf.WriteString("var countryRegionRange = [256]regionRange{\n")
	type pair struct {
		id  uint8
		off uint16
		n   uint16
	}
	var pairs []pair
	seen := map[uint8]bool{}
	for cc, list := range byCC {
		id, ok := ccToID[cc]
		if !ok || id == 0 || seen[id] {
			continue
		}
		seen[id] = true
		pairs = append(pairs, pair{id: id, off: list[0].id, n: uint16(len(list))})
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].id < pairs[j].id })
	for _, p := range pairs {
		fmt.Fprintf(buf, "\t%d: {off: %d, n: %d},\n", p.id, p.off, p.n)
	}
	buf.WriteString("}\n\n")
}

func bytesLit(b []byte) string {
	parts := make([]string, len(b))
	for i, v := range b {
		parts[i] = fmt.Sprintf("0x%02x", v)
	}
	return strings.Join(parts, ", ")
}

func fmtFloat(v float32) string {
	if v == 0 {
		return "0"
	}
	return strconv.FormatFloat(float64(v), 'g', -1, 32)
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "genregions: %v\n", err)
	os.Exit(1)
}
