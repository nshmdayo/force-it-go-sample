package main

import (
	"time"
)

// Particle replaces Perticle
type Particle struct {
	// Replacing Sphere inheritance with explicit coordinates
	X, Y, Z float64
	Radius  float64

	position *Vec
	velocity *Vec
	diff     *Vec
	t        int64 // timestamp in milliseconds
	deleteFlag bool

	// Constant 'a' from original code: 1/(4*Main.d*(Main.d+1))
	// Main.d = 10 -> a = 1/(4*10*11) = 1/440
	// We will calculate 'a' in the constructor or pass it
}

func NewParticle(x, y, z float64) *Particle {
	p := &Particle{
		X: x, Y: y, Z: z,
		position: NewVec(x, y, z),
		velocity: NewVec(0, 0, 0),
		diff:     NewVec(0, 0, 0),
		t:        time.Now().UnixNano() / int64(time.Millisecond),
	}
	return p
}

func (p *Particle) Move() {
	p.X += p.velocity.X
	p.Y += p.velocity.Y
	p.Z += p.velocity.Z
	p.diff.Reset()
}

func (p *Particle) Gravity() {
	if p.Z < 0 {
		p.X = p.position.X
		p.Y = p.position.Y
		p.Z = p.position.Z
		p.velocity.Reset()
	} else {
		// Vec a = new Vec(getTranslateX(), ...)
		currentPos := NewVec(p.X, p.Y, p.Z)
		v := NewVec(0, 0, 0)
		v.Sub(p.position, currentPos)
		v.Mult(0.0005)
		p.velocity.Add(v)
	}
}

func (p *Particle) SetAroundModule(v *Vec) {
	p.diff.Add(v)
}

func (p *Particle) ModuleGravity() {
	// a = 1/(4*d*(d+1)) where d=2 => a = 1/(4*2*3) = 1/24
	d := 2.0
	a := 1.0 / (4.0 * d * (d + 1.0))

	p.diff.Mult(a)
	p.velocity.Add(p.diff)
}

func (p *Particle) GetPosition() *Vec {
	return p.position
}

func (p *Particle) AddVelocity(v *Vec) {
	p.velocity.Add(v)
}

func (p *Particle) GetVelocity() *Vec {
	return p.velocity
}
