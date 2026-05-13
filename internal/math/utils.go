package math

import (
	"math"
	"math/rand"

	"github.com/infopek/agorio/internal/types"
)

func Sqrt(x types.Real) types.Real {
	return types.Real(math.Sqrt(float64(x)))
}

func Max(lhs, rhs types.Real) types.Real {
	return types.Real(math.Max(float64(lhs), float64(rhs)))
}

func RandRange(lo, hi int) int {
	return rand.Intn(hi-lo) + lo
}
