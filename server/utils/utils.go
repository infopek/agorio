package utils

import (
	"math/rand"
)

func RandIntRange(lo, hi int64) int64 {
	return rand.Int63n(hi-lo) + lo
}

func RandFloatRange(lo, hi float64) float64 {
	return rand.Float64()*(hi-lo) + lo
}

func GetRandomColor() [3]uint8 {
	return [3]uint8{
		uint8(RandIntRange(0, 256)),
		uint8(RandIntRange(0, 256)),
		uint8(RandIntRange(0, 256)),
	}
}
