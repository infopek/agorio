package physics

import (
	"github.com/infopek/agorio/internal/types"
)

type ShapeType types.Integer

const (
	ShapeTypeCircle ShapeType = iota
)

type Shape interface {
	GetType() ShapeType
	GetRadius() types.Real
	GetBody() *Body

	SetBody(b *Body)

	Initialize()
	ComputeMass(density types.Real)
}
