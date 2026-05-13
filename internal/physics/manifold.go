package physics

import (
	"github.com/infopek/agorio/internal/constants"
	"github.com/infopek/agorio/internal/math"
	"github.com/infopek/agorio/internal/types"
)

type Manifold struct {
	BodyA *Body
	BodyB *Body

	Penetration types.Real

	Normal math.Vector2
}

func NewManifold(bodyA, bodyB *Body) *Manifold {
	return &Manifold{
		BodyA: bodyA,
		BodyB: bodyB,
	}
}

func (m *Manifold) SetCollision() bool {
	// Vector from A to B
	delta := math.Sub(m.BodyB.Position, m.BodyA.Position)
	distSq := delta.LengthSq()
	radiusSum := m.BodyA.Shape.GetRadius() + m.BodyB.Shape.GetRadius()

	if distSq >= radiusSum*radiusSum {
		return false // not colliding
	}

	dist := math.Sqrt(distSq)
	if dist == 0 {
		// Circles on top of each other
		m.Normal = math.Vector2{X: 1.0, Y: 0.0}
		m.Penetration = m.BodyA.Shape.GetRadius()
	} else {
		m.Normal = math.Div(delta, dist)
		m.Penetration = radiusSum - dist
	}

	return true
}

func (m *Manifold) Resolve() {
	invMassSum := m.BodyA.InvMass + m.BodyB.InvMass
	if invMassSum == 0.0 {
		return // static objects
	}

	// Relative velocity along normal
	rv := math.Sub(m.BodyB.Velocity, m.BodyA.Velocity)
	contactVel := rv.Dot(m.Normal)
	if contactVel > 0.0 {
		return // not colliding
	}

	// Impulse scalar
	j := -contactVel / invMassSum
	impulse := math.Mul(m.Normal, j)

	// Apply impulse
	m.BodyA.Velocity = math.Sub(
		m.BodyA.Velocity,
		math.Mul(impulse, m.BodyA.InvMass),
	)
	m.BodyB.Velocity = math.Sub(
		m.BodyB.Velocity,
		math.Mul(impulse, m.BodyB.InvMass),
	)

	// Positional correction
	m.positionalCorrection()
}

func (m *Manifold) positionalCorrection() {
	correction := math.Max(m.Penetration-constants.PosCorrectionSlop, 0.0) /
		(m.BodyA.InvMass + m.BodyB.InvMass) * constants.PosCorrectionPercent
	correctionVec := math.Mul(m.Normal, correction)

	m.BodyA.Position = math.Add(m.BodyA.Position, math.Mul(correctionVec, m.BodyA.InvMass))
	m.BodyB.Position = math.Add(m.BodyB.Position, math.Mul(correctionVec, m.BodyB.InvMass))
}
