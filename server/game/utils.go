package game

import (
	"math"
	"math/rand"
)

func RandIntRange(lo, hi int64) int64 {
	return rand.Int63n(hi-lo) + lo
}

func RandFloatRange(lo, hi float64) float64 {
	return rand.Float64()*(hi-lo) + lo
}

func Clamp(val, lo, hi float64) float64 {
	if val < lo {
		return lo
	}
	if val > hi {
		return hi
	}
	return val
}

/** Radius
 *
 * Universal radius calculation for every physical entity
 *
 */
func Radius(mass float64) float64 {
	return RadiusScale * math.Sqrt(mass)
}
