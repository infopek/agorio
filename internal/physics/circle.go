package physics

import (
	"math"

	"github.com/infopek/agorio/internal/types"
)

// Implements Shape
type Circle struct {
	Radius types.Real
	Body   *Body // is owned by this Body
}

func NewCircle(radius types.Real) Circle {
	return Circle{
		Radius: radius,
		Body:   nil,
	}
}

func (c *Circle) GetType() ShapeType {
	return ShapeTypeCircle
}

func (c *Circle) GetRadius() types.Real {
	return c.Radius
}

func (c *Circle) GetBody() *Body {
	return c.Body
}

func (c *Circle) SetBody(b *Body) {
	c.Body = b
}

func (c *Circle) Initialize() {
	c.ComputeMass(1.0)
}

func (c *Circle) ComputeMass(density types.Real) {
	c.Body.Mass = math.Pi * c.Radius * c.Radius * density
	c.Body.InvMass = 0.0
	if c.Body.Mass != 0.0 {
		c.Body.InvMass = 1.0 / c.Body.Mass
	}
}
