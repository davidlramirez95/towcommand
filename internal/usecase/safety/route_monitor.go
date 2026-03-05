package safetyuc

import "math"

// DefaultDeviationThresholdKm is the default distance threshold for route deviation detection.
const DefaultDeviationThresholdKm = 2.0

// RoutePoint represents a geographic coordinate used for route monitoring.
type RoutePoint struct {
	Lat float64
	Lng float64
}

// CheckRouteDeviation checks if the current position deviates from the expected
// corridor between pickup and dropoff by more than thresholdKm.
//
// The corridor is defined by the line segment from pickup to dropoff. The
// function computes the shortest distance from the current point to this
// segment and returns true if that distance exceeds the threshold.
func CheckRouteDeviation(current, pickup, dropoff RoutePoint, thresholdKm float64) bool {
	dist := pointToSegmentDistanceKm(current, pickup, dropoff)
	return dist > thresholdKm
}

// pointToSegmentDistanceKm computes the shortest distance in km from a point
// to the line segment defined by two endpoints, using haversine-based
// calculations.
//
// Algorithm:
//  1. Project the point onto the infinite line through the segment endpoints.
//  2. If the projection parameter t is in [0,1], use the perpendicular
//     distance to the line.
//  3. Otherwise, use the distance to the nearest endpoint.
func pointToSegmentDistanceKm(p, a, b RoutePoint) float64 {
	// Convert to approximate Cartesian coordinates (km) centred on point a.
	// This is sufficiently accurate for Philippine latitudes and short distances.
	latMid := degreesToRad((a.Lat + b.Lat) / 2.0)
	kmPerDegLat := 111.32
	kmPerDegLng := 111.32 * math.Cos(latMid)

	ax, ay := 0.0, 0.0
	bx := (b.Lng - a.Lng) * kmPerDegLng
	by := (b.Lat - a.Lat) * kmPerDegLat
	px := (p.Lng - a.Lng) * kmPerDegLng
	py := (p.Lat - a.Lat) * kmPerDegLat

	// Vector AB and AP.
	abx := bx - ax
	aby := by - ay
	apx := px - ax
	apy := py - ay

	ab2 := abx*abx + aby*aby
	if ab2 == 0 {
		// Segment is a single point.
		return math.Sqrt(apx*apx + apy*apy)
	}

	// Parameter t of the projection of P onto AB.
	t := (apx*abx + apy*aby) / ab2

	if t < 0 {
		// Closest to endpoint A.
		return math.Sqrt(apx*apx + apy*apy)
	}
	if t > 1 {
		// Closest to endpoint B.
		dx := px - bx
		dy := py - by
		return math.Sqrt(dx*dx + dy*dy)
	}

	// Perpendicular distance.
	projX := ax + t*abx
	projY := ay + t*aby
	dx := px - projX
	dy := py - projY
	return math.Sqrt(dx*dx + dy*dy)
}

// degreesToRad converts degrees to radians.
func degreesToRad(deg float64) float64 {
	return deg * math.Pi / 180.0
}
