package utils

import "math"

type Point struct {
	Lat float64
	Lon float64
}

func Haversine(a, b Point) float64 {
	const R = 6371.0
	const degToRad = math.Pi / 180.0

	dLat := (b.Lat - a.Lat) * degToRad
	dLon := (b.Lon - a.Lon) * degToRad
	lat1 := a.Lat * degToRad
	lat2 := b.Lat * degToRad

	sinDLat := math.Sin(dLat / 2)
	sinDLon := math.Sin(dLon / 2)

	h := sinDLat*sinDLat + math.Cos(lat1)*math.Cos(lat2)*sinDLon*sinDLon
	return 2 * R * math.Asin(math.Sqrt(h))
}

func NearestNeighborTSP(startIdx int, points []Point) float64 {
	n := len(points)
	if n < 2 {
		return 0
	}

	visited := make([]bool, n)
	current := startIdx
	totalDist := 0.0

	for step := 0; step < n-2; step++ {
		visited[current] = true
		minDist := math.MaxFloat64
		next := -1

		for j := 0; j < n-1; j++ {
			if !visited[j] {
				dist := Haversine(points[current], points[j])
				if dist < minDist {
					minDist = dist
					next = j
				}
			}
		}

		if next == -1 {
			break
		}
		totalDist += minDist
		current = next
	}

	totalDist += Haversine(points[current], points[n-1])
	return totalDist
}
