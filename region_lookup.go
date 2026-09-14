package gogeo

import "bytes"

type regionRange struct {
	off uint16
	n   uint16
}

// Regions is the public list of all subdivisions. Index equals Region.ID.
var Regions = regions[:]

// RegionByID returns the region with the given ID, or the undefined region.
func RegionByID(id uint16) *Region {
	if int(id) >= len(regions) {
		return &regions[0]
	}
	return &regions[id]
}

// RegionByCode returns the region for an ISO 3166-2 code (for example "US-CA").
// "GB-*" and "UK-*" resolve to the same United Kingdom subdivisions.
// Unknown codes resolve to the undefined region (never nil).
func RegionByCode(code string) *Region {
	b0, b1, suf, ok := parseRegionCode(code)
	if !ok {
		return &regions[0]
	}
	if b0 < 'A' || b0 > 'Z' || b1 < 'A' || b1 > 'Z' {
		return &regions[0]
	}
	return lookupRegionSuffix(iso2RegionRange[b0-'A'][b1-'A'], suf)
}

func parseRegionCode(code string) (b0, b1 byte, suf string, ok bool) {
	if len(code) < 4 || len(code) > 6 || code[2] != '-' {
		return 0, 0, "", false
	}
	b0, b1 = code[0], code[1]
	suf = code[3:]
	if len(suf) == 0 || len(suf) > 3 {
		return 0, 0, "", false
	}
	return b0, b1, suf, true
}

func lookupRegionSuffix(rng regionRange, suf string) *Region {
	if rng.n == 0 {
		return &regions[0]
	}
	var buf [3]byte
	key := buf[:copy(buf[:], suf)]
	lo, hi := int(rng.off), int(rng.off)+int(rng.n)
	for lo < hi {
		mid := lo + (hi-lo)/2
		r := &regions[mid]
		cmp := bytes.Compare(r.suffix[:r.sufN], key)
		if cmp == 0 {
			return r
		}
		if cmp < 0 {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return &regions[0]
}
