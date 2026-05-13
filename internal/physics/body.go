package physics

import (
	"github.com/infopek/agorio/internal/math"
	"github.com/infopek/agorio/internal/types"
)

type Body struct {
	Position math.Vector2 `json:"position"`
	Velocity math.Vector2 `json:"-"`
	Force    math.Vector2 `json:"-"`
	Mass     types.Real   `json:"mass"`
	InvMass  types.Real   `json:"-"`
	Shape    Shape        `json:"-"` // owns a Shape
}

func NewBody(shape Shape, position math.Vector2) *Body {
	body := Body{
		Position: position,
		Velocity: math.Vector2{X: 0.0, Y: 0.0},
		Force:    math.Vector2{X: 0.0, Y: 0.0},
		Shape:    shape,
	}

	body.Shape.SetBody(&body)
	body.Shape.Initialize()

	return &body
}

func (b *Body) ApplyForce(f math.Vector2) {
	b.Force.Addi(f)
}

func (b *Body) ApplyImpulse(impulse math.Vector2, contactVector math.Vector2) {
	b.Velocity.Addsi(impulse, b.InvMass)
}
