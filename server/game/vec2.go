package game

import (
	"math"
)

type Vec2 struct {
	X float64
	Y float64
}

func (v Vec2) Add(other Vec2) Vec2 {
	return Vec2{
		X: v.X + other.X,
		Y: v.Y + other.Y,
	}
}

func (v Vec2) Sub(other Vec2) Vec2 {
	return Vec2{
		X: v.X - other.X,
		Y: v.Y - other.Y,
	}
}

func (v Vec2) Scale(s float64) Vec2 {
	return Vec2{
		X: v.X * s,
		Y: v.Y * s,
	}
}

func (v Vec2) Normalize() Vec2 {
	mag := v.Magnitude()
	if mag == 0.0 {
		return Vec2{}
	}

	return Vec2{
		X: v.X / mag,
		Y: v.Y / mag,
	}
}

func (v Vec2) Magnitude() float64 {
	return math.Sqrt(v.MagnitudeSq())
}

func (v Vec2) MagnitudeSq() float64 {
	return v.X*v.X + v.Y*v.Y
}

func (v Vec2) DistanceTo(other Vec2) float64 {
	temp := other.Sub(v)
	return temp.Magnitude()
}

func (v Vec2) Lerp(other Vec2, t float64) Vec2 {
	return Vec2{
		X: v.X * t + other.X * (1.0 - t),
		Y: v.Y * t + other.Y * (1.0 - t),
	}
}
