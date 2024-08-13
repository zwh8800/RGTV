package util

import (
	"math"
)

// CubicBezier calculates the cubic Bezier curve value for a given t
func CubicBezier(t, p0, p1, p2, p3 float64) float64 {
	return math.Pow(1-t, 3)*p0 +
		3*math.Pow(1-t, 2)*t*p1 +
		3*(1-t)*math.Pow(t, 2)*p2 +
		math.Pow(t, 3)*p3
}

// EaseInOut calculates the ease-in-out value for a given t
func EaseInOut(t float64) float64 {
	// Control points for ease-in-out
	p0, p1, p2, p3 := 0.0, 0.42, 0.58, 1.0
	return CubicBezier(t, p0, p1, p2, p3)
}

// EaseOut calculates the ease-out value for a given t
func EaseOut(t float64) float64 {
	// Control points for ease-out
	p0, p1, p2, p3 := 0.0, 0.0, 0.58, 1.0
	return CubicBezier(t, p0, p1, p2, p3)
}

// EaseIn calculates the ease-in value for a given t
func EaseIn(t float64) float64 {
	// Control points for ease-in
	p0, p1, p2, p3 := 0.0, 0.42, 1.0, 1.0
	return CubicBezier(t, p0, p1, p2, p3)
}
