// Package utcoffset resolves a stable standard-time UTC offset for an IANA zone.
package utcoffset

import (
	"math"
	"time"
	_ "time/tzdata"
)

// StandardOffset returns minutes east of UTC for standard time (not DST).
// Unknown zone names fall back to round(lon/15)*60.
func StandardOffset(zone string, lon float64) int16 {
	if zone != "" {
		if off, ok := cache[zone]; ok {
			return off
		}
		if loc, err := time.LoadLocation(zone); err == nil {
			off := standardMinutes(loc)
			cache[zone] = off
			return off
		}
	}
	return int16(math.Round(lon/15) * 60)
}

var cache = map[string]int16{}

func standardMinutes(loc *time.Location) int16 {
	fallback := 0
	for month := time.January; month <= time.December; month++ {
		t := time.Date(2026, month, 15, 12, 0, 0, 0, loc)
		_, sec := t.Zone()
		if !t.IsDST() {
			return int16(sec / 60)
		}
		if month == time.January {
			fallback = sec / 60
		}
	}
	return int16(fallback)
}
