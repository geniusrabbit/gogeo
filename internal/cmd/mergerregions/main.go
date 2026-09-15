// Command mergerregions merges Debian iso-codes and oodavid/iso-3166-2 into data/regions.json.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

type debianFile struct {
	Items []debianRegion `json:"3166-2"`
}

type debianRegion struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Parent string `json:"parent"`
}

type oodavidRegion struct {
	Code string   `json:"code"`
	Name string   `json:"name"`
	Lat  *float64 `json:"lat"`
	Lng  *float64 `json:"lng"`
}

type sourceMeta struct {
	ID      string `json:"id"`
	URL     string `json:"url"`
	License string `json:"license,omitempty"`
	Field   string `json:"field,omitempty"`
}

type mergedFile struct {
	Sources []sourceMeta   `json:"sources"`
	Regions []mergedRegion `json:"regions"`
}

type mergedRegion struct {
	Code   string            `json:"code"`
	Name   string            `json:"name"`
	Names  map[string]string `json:"names,omitempty"`
	Type   string            `json:"type,omitempty"`
	Parent string            `json:"parent,omitempty"`
	Lat    *float64          `json:"lat,omitempty"`
	Lng    *float64          `json:"lng,omitempty"`
}

func main() {
	debian := flag.String("debian", "data/iso_3166-2.json", "Debian iso-codes JSON")
	oodavid := flag.String("oodavid", "data/oodavid_iso_3166_2.js", "oodavid iso_3166_2.js or JSON")
	out := flag.String("out", "data/regions.json", "merged JSON")
	flag.Parse()

	items, err := loadDebian(*debian)
	if err != nil {
		fatal(err)
	}
	extra, err := loadOodavid(*oodavid)
	if err != nil {
		fatal(err)
	}

	merged := mergedFile{
		Sources: []sourceMeta{
			{ID: "iso-codes", URL: "https://salsa.debian.org/iso-codes-team/iso-codes", License: "LGPL-2.1-or-later"},
			{ID: "oodavid", URL: "https://github.com/oodavid/iso-3166-2", Field: "lat,lng,name"},
		},
		Regions: make([]mergedRegion, 0, len(items)),
	}
	for _, item := range items {
		merged.Regions = append(merged.Regions, overlay(item, extra[item.Code]))
	}

	raw, err := json.MarshalIndent(merged, "", "  ")
	if err != nil {
		fatal(err)
	}
	raw = append(raw, '\n')
	if err := os.WriteFile(*out, raw, 0644); err != nil {
		fatal(err)
	}
}

func loadDebian(path string) ([]debianRegion, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var file debianFile
	if err := json.Unmarshal(raw, &file); err != nil {
		return nil, err
	}
	if len(file.Items) == 0 {
		return nil, fmt.Errorf("no debian regions in %s", path)
	}
	return file.Items, nil
}

func loadOodavid(path string) (map[string]oodavidRegion, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	js, err := jsObject(raw)
	if err != nil {
		return nil, err
	}
	var obj map[string]oodavidRegion
	if err := json.Unmarshal(js, &obj); err != nil {
		return nil, fmt.Errorf("oodavid: %w", err)
	}
	out := make(map[string]oodavidRegion, len(obj))
	for code, rec := range obj {
		if !strings.Contains(code, "-") {
			continue
		}
		if rec.Code == "" {
			rec.Code = code
		}
		out[code] = rec
	}
	return out, nil
}

func jsObject(raw []byte) ([]byte, error) {
	s := string(raw)
	if i := strings.Index(s, "{"); i >= 0 {
		s = s[i:]
	}
	var b strings.Builder
	inStr := false
	esc := false
	i := 0
	for i < len(s) {
		c := s[i]
		if inStr {
			b.WriteByte(c)
			if esc {
				esc = false
			} else if c == '\\' {
				esc = true
			} else if c == '"' {
				inStr = false
			}
			i++
			continue
		}
		if c == '/' && i+1 < len(s) && s[i+1] == '/' {
			if j := strings.IndexByte(s[i:], '\n'); j >= 0 {
				i += j + 1
			} else {
				break
			}
			continue
		}
		if c == '"' {
			inStr = true
			b.WriteByte(c)
			i++
			continue
		}
		b.WriteByte(c)
		i++
	}
	out := []byte(strings.TrimSpace(b.String()))
	out = bytes.TrimRight(out, "; \n\t\r")
	if len(out) == 0 || out[0] != '{' {
		return nil, fmt.Errorf("oodavid: no JSON object")
	}
	return out, nil
}

func overlay(item debianRegion, extra oodavidRegion) mergedRegion {
	out := mergedRegion{
		Code:   item.Code,
		Name:   item.Name,
		Type:   item.Type,
		Parent: item.Parent,
		Names:  map[string]string{},
	}
	if out.Name != "" {
		out.Names["en"] = out.Name
	}
	if extra.Name != "" {
		for _, alt := range nameVariants(extra.Name) {
			if alt == "" || alt == out.Name {
				continue
			}
			if _, ok := out.Names["alt"]; !ok {
				out.Names["alt"] = alt
			} else if out.Names["alt"] != alt {
				out.Names["alt2"] = alt
			}
		}
		if out.Name == "" {
			out.Name = extra.Name
			out.Names["en"] = extra.Name
		}
	}
	if extra.Lat != nil && extra.Lng != nil && (*extra.Lat != 0 || *extra.Lng != 0) {
		out.Lat = extra.Lat
		out.Lng = extra.Lng
	}
	if len(out.Names) == 0 {
		out.Names = nil
	}
	return out
}

func nameVariants(name string) []string {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil
	}
	out := []string{name}
	if i := strings.IndexByte(name, '['); i >= 0 {
		j := strings.LastIndexByte(name, ']')
		if j > i+1 {
			inner := strings.TrimSpace(name[i+1 : j])
			outer := strings.TrimSpace(name[:i])
			if inner != "" {
				out = append(out, inner)
			}
			if outer != "" {
				out = append(out, outer)
			}
		}
	}
	return out
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "mergerregions: %v\n", err)
	os.Exit(1)
}
