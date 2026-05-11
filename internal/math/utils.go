package math

import (
	"math/rand"
)

func RandRange(lo, hi int) int {
	return rand.Intn(hi-lo) + lo
}
