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

// HaversineKm menghitung jarak dua titik (kilometer)
func HaversineKm(a, b Point) float64 {
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

// NearestNeighborFromUser menjalankan algoritma greedy Nearest Neighbor mulai dari titik user.
func NearestNeighborFromUser(user Point, merchants []Point) []int {
	n := len(merchants)
	if n == 0 {
		return nil
	}

	visited := make([]bool, n)
	order := make([]int, 0, n)
	current := user

	for len(order) < n {
		minIdx := -1
		minDist := math.MaxFloat64
		for i := 0; i < n; i++ {
			if visited[i] {
				continue
			}
			dist := HaversineKm(current, merchants[i])
			if dist < minDist {
				minDist = dist
				minIdx = i
			}
		}
		if minIdx == -1 {
			break
		}
		visited[minIdx] = true
		order = append(order, minIdx)
		current = merchants[minIdx]
	}
	return order
}
