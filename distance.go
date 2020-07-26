package main

import (
	"math"
	"unicode/utf8"
)

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

func distanceInKiloMeters(lat1, lon1, lat2, lon2 float64) float64 {
	const r = float64(6378.1)

	lat1 = lat1 * math.Pi / 180
	lon1 = lon1 * math.Pi / 180
	lat2 = lat2 * math.Pi / 180
	lon2 = lon2 * math.Pi / 180

	h := hsin(lat2-lat1) + math.Cos(lat1)*math.Cos(lat2)*hsin(lon2-lon1)

	return 2 * r * math.Asin(math.Sqrt(h))
}

func editDistance(a, b string) int {
	var arr [64]int
	rc := utf8.RuneCountInString(b) + 1
	if rc > len(arr) {
		return 32
	}
	f := arr[:rc]

	for j := range f {
		f[j] = j
	}

	for _, ca := range a {
		j := 1
		fj1 := f[0]
		f[0]++
		for _, cb := range b {
			mn := min(f[j]+1, f[j-1]+1)
			if cb != ca {
				mn = min(mn, fj1+1)
			} else {
				mn = min(mn, fj1)
			}

			fj1, f[j] = f[j], mn
			j++
		}
	}

	return f[len(f)-1]
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
