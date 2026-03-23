package math

import (
	"math/rand"
)

type Vector2 struct {
	X int
	Y int
}

func RandRange(lo, hi int) int {
	return rand.Intn(hi - lo) + lo
}
