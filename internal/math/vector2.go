package math

import (
	"github.com/infopek/agorio/internal/types"
)

type Vector2 struct {
	X types.Real `json:"x"`
	Y types.Real `json:"y"`
}

/**
 * Returns the squared length of v
 */
func (v Vector2) LengthSq() types.Real {
	return v.X*v.X + v.Y*v.Y
}

/**
 * Returns a new vector that is the addition of lhs and rhs
 */
func Add(lhs, rhs Vector2) Vector2 {
	return Vector2{
		X: lhs.X + rhs.X,
		Y: lhs.Y + rhs.Y,
	}
}

/**
 * Returns a new vector that is the addition of lhs and rhs
 */
func Sub(lhs, rhs Vector2) Vector2 {
	return Vector2{
		X: lhs.X - rhs.X,
		Y: lhs.Y - rhs.Y,
	}
}

/**
 * Returns a new vector that's elements are v's elements
 *  multiplied by scalar
 */
func Mul(v Vector2, scalar types.Real) Vector2 {
	return Vector2{
		X: v.X * scalar,
		Y: v.Y * scalar,
	}
}

/**
 * Adds other to v, and returns the result
 */
func (v *Vector2) Addi(other Vector2) *Vector2 {
	v.X += other.X
	v.Y += other.Y
	return v
}

/**
 * Multiplies v's elements by scalar, and returns the result
 */
func (v *Vector2) Mul(scalar types.Real) *Vector2 {
	v.X *= scalar
	v.Y *= scalar
	return v
}

/**
 * Adds other * s to v
 */
func (v *Vector2) Addsi(other Vector2, scalar types.Real) *Vector2 {
	v.X += other.X * scalar
	v.Y += other.Y * scalar
	return v
}

/**
 * Returns a new vector with elements divided by scalar
 */
func Div(v Vector2, scalar types.Real) Vector2 {
	return Vector2{
		X: v.X / scalar,
		Y: v.Y / scalar,
	}
}

/**
 * Returns the dot product of v and other
 */
func (v *Vector2) Dot(other Vector2) types.Real {
	return v.X*other.X + v.Y*other.Y
}

/**
 * Returns the normalized vector of v
 */
func (v Vector2) Normalized() Vector2 {
	length := Sqrt(v.LengthSq())
	return Vector2{
		X: v.X / length,
		Y: v.Y / length,
	}
}

/**
 * Returns a new vector that has randomized elements
 *  that are in the range [loX, hiX), and [loY, hiY)
 */
func RandVector2(loX, hiX, loY, hiY int) Vector2 {
	return Vector2{
		X: types.Real(RandRange(loX, hiX)),
		Y: types.Real(RandRange(loY, hiY)),
	}
}
