package main

import (
	"testing"
)

func TestVecAdd(t *testing.T) {
	v1 := NewVec(1, 2, 3)
	v2 := NewVec(4, 5, 6)
	v1.Add(v2)
	if v1.X != 5 || v1.Y != 7 || v1.Z != 9 {
		t.Errorf("Vec.Add failed: got %v, expected {5, 7, 9}", v1)
	}
}

func TestVecMag(t *testing.T) {
	v := NewVec(3, 4, 0)
	if v.Mag() != 25 {
		t.Errorf("Vec.Mag failed: got %f, expected 25", v.Mag())
	}
	if v.Dist() != 5 {
		t.Errorf("Vec.Dist failed: got %f, expected 5", v.Dist())
	}
}

func TestParticleMove(t *testing.T) {
	p := NewParticle(0, 0, 0)
	p.AddVelocity(NewVec(1, 1, 1))
	p.Move()
	if p.X != 1 || p.Y != 1 || p.Z != 1 {
		t.Errorf("Particle.Move failed: got {%f, %f, %f}", p.X, p.Y, p.Z)
	}
}

func TestParticleGravity(t *testing.T) {
	// If Z < 0, it should reset
	p := NewParticle(10, 10, 10)
	p.Z = -5
	p.Gravity()
	if p.Z != 10 {
		t.Errorf("Particle.Gravity reset failed: got Z=%f, expected 10", p.Z)
	}
}
