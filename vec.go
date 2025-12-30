package main

import "math"

type Vec struct {
	X, Y, Z float64
}

func NewVec(x, y, z float64) *Vec {
	return &Vec{X: x, Y: y, Z: z}
}

func (v *Vec) Add(target *Vec) {
	v.X += target.X
	v.Y += target.Y
	v.Z += target.Z
}

func (v *Vec) Sub(target, current *Vec) {
	v.X = target.X - current.X
	v.Y = target.Y - current.Y
	v.Z = target.Z - current.Z
}

func (v *Vec) Mult(k float64) {
	v.X *= k
	v.Y *= k
	v.Z *= k
}

func (v *Vec) Dot(another *Vec) float64 {
	return v.X*another.X + v.Y*another.Y + v.Z*another.Z
}

func (v *Vec) Mag() float64 {
	return math.Pow(v.X, 2) + math.Pow(v.Y, 2) + math.Pow(v.Z, 2)
}

func (v *Vec) Dist() float64 {
	return math.Sqrt(v.Mag())
}

func (v *Vec) Normalize() {
	d := v.Dist()
	if d != 0 {
		v.X /= d
		v.Y /= d
		v.Z /= d
	}
}

func (v *Vec) Copy(another *Vec) {
	v.X = another.X
	v.Y = another.Y
	v.Z = another.Z
}

func (v *Vec) Reset() {
	v.X = 0
	v.Y = 0
	v.Z = 0
}
