package math

import (
	"math/rand"

	"github.com/infopek/agorio/internal/types"
)

type Vector2 struct {
	X types.Real
	Y types.Real
}

func RandRange(lo, hi int) int {
	return rand.Intn(hi-lo) + lo
}

func RandVector2(loX, hiX, loY, hiY int) Vector2 {
	return Vector2{
		X: types.Real(RandRange(loX, hiX)),
		Y: types.Real(RandRange(loY, hiY)),
	}
}
