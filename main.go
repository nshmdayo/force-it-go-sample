package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// Constants
const (
	SCENE_WIDTH  = 1920
	SCENE_HEIGHT = 1080
	CAMERA_X     = 0
	CAMERA_Y     = -100
	CAMERA_Z     = -250

	makeForceNum     = 50
	FORCE_POWER_RATE = 0.4
	ballradius       = 5.0
	forceradius      = 1.0
	s_const          = (ballradius + forceradius) * (ballradius + forceradius)
	d_const          = 2

	m1 = 1.0
	m2 = 1.0
	a_mass = 1.0 / (m1 + m2)
)

// Affine transform struct (simplified for 3D points)
type Affine struct {
	m00, m01, m02, tx float64
	m10, m11, m12, ty float64
	m20, m21, m22, tz float64
}

func NewAffineIdentity() *Affine {
	return &Affine{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
	}
}

// Rotate X
func NewRotateX(theta float64) *Affine {
	c := math.Cos(theta)
	s := math.Sin(theta)
	return &Affine{
		1, 0, 0, 0,
		0, c, -s, 0,
		0, s, c, 0,
	}
}

// Scale
func NewScale(sx, sy, sz float64) *Affine {
	return &Affine{
		sx, 0, 0, 0,
		0, sy, 0, 0,
		0, 0, sz, 0,
	}
}

// Translate
func NewTranslate(tx, ty, tz float64) *Affine {
	return &Affine{
		1, 0, 0, tx,
		0, 1, 0, ty,
		0, 0, 1, tz,
	}
}

func (a *Affine) Transform(x, y, z float64) (float64, float64, float64) {
	nx := a.m00*x + a.m01*y + a.m02*z + a.tx
	ny := a.m10*x + a.m11*y + a.m12*z + a.ty
	nz := a.m20*x + a.m21*y + a.m22*z + a.tz
	return nx, ny, nz
}

type Simulation struct {
	wallBalls []*Particle
	forces    []*Particle
	kinect    *KinectDevice

	rightVec   *Vec
	right      *Vec
	left       *Vec
	spineS     *Vec
	base       *Vec
	shoulderR  *Vec

	rotate    *Affine
	translate *Affine
	scale     *Affine
}

func NewSimulation() *Simulation {
	sim := &Simulation{
		wallBalls: make([]*Particle, 0),
		forces:    make([]*Particle, 0),
		kinect:    NewKinectDevice(),
		rightVec:  NewVec(0, 0, 0),
		right:     NewVec(0, 0, 0),
		left:      NewVec(0, 0, 0),
		spineS:    NewVec(0, 0, 0),
		base:      NewVec(0, 0, 0),
		shoulderR: NewVec(0, 0, 0),
		rotate:    NewAffineIdentity(),
		translate: NewAffineIdentity(),
		scale:     NewAffineIdentity(),
	}

	// Initialize WallBalls
	for y := -20; y < 0; y++ {
		for x := -20; x < 20; x++ {
			p := NewParticle(float64(x)*ballradius*2, float64(y)*ballradius*2, 0)
			p.Radius = ballradius
			sim.wallBalls = append(sim.wallBalls, p)
		}
	}
	return sim
}

func (sim *Simulation) Calibration() {
	if !sim.kinect.Update() {
		return
	}
	// bX, bY, bZ, _ := sim.kinect.GetJoint(JointType_SpineBase)
	_, _, bZ, _ := sim.kinect.GetJoint(JointType_SpineBase)
	mX, mY, mZ, _ := sim.kinect.GetJoint(JointType_SpineMid)

	tmp := bZ / 2.0
	if tmp > 1.0 {
		tmp = 1.0
	}
	theta := math.Acos(tmp)
	// System.out.println(Math.toDegrees(theta));
	fmt.Printf("Calibration Theta: %f\n", theta*180/math.Pi)

	sim.rotate = NewRotateX(theta)

	// m2 = rotate.transform(m2)
	_, my2, _ := sim.rotate.Transform(mX, mY, mZ)

	// translate = new Affine(1, 0, 0, 0, 0, 1, 0, -m.getY(), 0, 0, 1, 0);
	// The original Java code says: new Affine(1, 0, 0, 0, 0, 1, 0, -m.getY(), 0, 0, 1, 0)
	// Which corresponds to Translate(0, -m.getY(), 0) if it is row-major/column-major?
	// JavaFX Affine(mxx, mxy, mxz, tx, myx, myy, myz, ty, mzx, mzy, mzz, tz)
	// So tx=0, ty=-m.getY(), tz=0. But wait, in the calibration lambda:
	// translate = new Affine(1, 0, 0, 0, 0, 1, 0, -m.getY(), 0, 0, 1, 0);
	// m.getY() is the raw Y coordinate.
	// But later they use it.
	// NOTE: m2 is transformed by rotate. The translate uses -m.getY() (raw m).

	sim.translate = NewTranslate(0, -my2, 0) // Using transformed Y? or raw? Java code used `m.getY()` but m was raw.
	// Wait, code:
	// Point3D m2 = new Point3D(m.getX(), m.getY(), m.getZ());
	// m2 = rotate.transform(m2);
	// translate = new Affine(..., -m.getY(), ...);
	// It seems to use raw m.getY(). But usually you translate relative to the transformed point.
	// Let's assume raw as in code.

	sim.scale = NewScale(100, -100, -100)
}

func (sim *Simulation) Recognize() {
	sim.kinect.Update()

	rx, ry, rz, rState := sim.kinect.GetJoint(JointType_HandRight)
	lx, ly, lz, _ := sim.kinect.GetJoint(JointType_HandLeft)
	ssx, ssy, ssz, _ := sim.kinect.GetJoint(JointType_SpineShoulder)
	bx, by, bz, _ := sim.kinect.GetJoint(JointType_SpineBase)
	srx, sry, srz, _ := sim.kinect.GetJoint(JointType_ShoulderRight)

	sim.right = NewVec(rx, ry, rz)
	sim.left = NewVec(lx, ly, lz)
	sim.spineS = NewVec(ssx, ssy, ssz)
	sim.base = NewVec(bx, by, bz)
	sim.shoulderR = NewVec(srx, sry, srz)

	sim.rightVec.Reset()

	if rState == HandState_Open {
		sim.rightVec.Sub(sim.right, sim.shoulderR)
	}
}

func (sim *Simulation) MakeForce() {
	rightHandX, rightHandY, rightHandZ := sim.TransformPoint(sim.right.X, sim.right.Y, sim.right.Z)
	leftHandX, leftHandY, leftHandZ := sim.TransformPoint(sim.left.X, sim.left.Y, sim.left.Z)
	// spineShoulderX, spineShoulderY, spineShoulderZ := sim.TransformPoint(sim.spineS.X, sim.spineS.Y, sim.spineS.Z)
	spineBaseX, spineBaseY, spineBaseZ := sim.TransformPoint(sim.base.X, sim.base.Y, sim.base.Z)

	sim.rightVec.Mult(FORCE_POWER_RATE)
	rightVX, rightVY, rightVZ := sim.TransformPoint(sim.rightVec.X, sim.rightVec.Y, sim.rightVec.Z)

	d := math.Sqrt(math.Pow(spineBaseX-leftHandX, 2) + math.Pow(spineBaseY-leftHandY, 2) + math.Pow(spineBaseZ-leftHandZ, 2))

	if rightVX != 0 && rightVY != 0 && rightVZ != 0 {
		for i := 0; i < makeForceNum; i++ {
			d2x := d * (rand.Float64() - 0.5) * 0.05
			d2y := d * (rand.Float64() - 0.5) * 0.05
			d2z := d * (rand.Float64() - 0.5) * 0.05

			f := NewParticle(rightHandX*d2x, rightHandY*d2y, rightHandZ*d2z)
			f.Radius = forceradius
			temp := NewVec(rightVX, rightVY, rightVZ)
			f.AddVelocity(temp)
			sim.forces = append(sim.forces, f)
		}
	}
}

func (sim *Simulation) TransformPoint(x, y, z float64) (float64, float64, float64) {
	x, y, z = sim.rotate.Transform(x, y, z)
	x, y, z = sim.translate.Transform(x, y, z)
	x, y, z = sim.scale.Transform(x, y, z)
	return x, y, z
}

func (sim *Simulation) Collision() {
	for i := 0; i < len(sim.wallBalls); i++ {
		for j := 0; j < len(sim.forces); j++ {
			if sim.CollisionFlag(sim.wallBalls[i], sim.forces[j]) {
				v1 := NewVec(0, 0, 0)
				v2 := NewVec(0, 0, 0)

				v1.Copy(sim.wallBalls[i].GetVelocity())
				v1.Mult(m1 - m2)

				v2.Copy(sim.forces[j].GetVelocity())
				v2.Mult(2 * m2)

				v1.Add(v2)
				v1.Mult(a_mass)

				sim.wallBalls[i].AddVelocity(v1)
				sim.forces[j].deleteFlag = true
			}
		}
	}

	// Remove flagged forces
	newForces := make([]*Particle, 0)
	for _, f := range sim.forces {
		if !f.deleteFlag {
			newForces = append(newForces, f)
		}
	}
	sim.forces = newForces
}

func (sim *Simulation) CollisionFlag(ball, force *Particle) bool {
	dx := math.Pow(ball.X-force.X, 2)
	dy := math.Pow(ball.Y-force.Y, 2)
	dz := math.Pow(ball.Z-force.Z, 2)
	d := dx + dy + dz
	return d < s_const
}

func (sim *Simulation) Move() {
	for _, p := range sim.forces {
		p.Move()
	}

	// Remove old forces (> 5000ms)
	now := time.Now().UnixNano() / int64(time.Millisecond)
	newForces := make([]*Particle, 0)
	for _, f := range sim.forces {
		if now-f.t <= 5000 {
			newForces = append(newForces, f)
		}
	}
	sim.forces = newForces
}

func (sim *Simulation) UpdateWallBalls() {
	// Interaction between wall balls
	for y := d_const; y < 20-d_const; y++ {
		for x := d_const; x < 40-d_const; x++ {
			// WallBalls grid is -20 to 0 (height 20) and -20 to 20 (width 40).
			// Size is 20*40 = 800.
			// Indices: y from 0 to 19, x from 0 to 39?
			// Init: y = -20 to -1. x = -20 to 19.
			// Index calculation in Java: int a = (y+j) * 40 + x+i;
			// The original loop for y is `int y = d; y < 20-d`.
			// The grid creation: y=-20..-1.
			// This index logic seems suspect in original code if not mapped correctly.
			// Let's assume 1D array layout row by row.
			// y index 0 corresponds to y=-20.

			// Let's map 2D to 1D properly.
			// Array size: 20 * 40 = 800.
			// index = (y_idx) * 40 + x_idx.
			// y_idx goes 0..19. x_idx goes 0..39.

			// The loop uses `y` and `x` variables but they seem to be indices into the array, not coordinates.
			// `for(int y = d; y < 20-d; y++)`
			// `for(int x = d; x < 40-d; x++)`
			// This matches the dimensions 20 and 40.

			for j := -d_const; j <= d_const; j++ {
				for i := -d_const; i <= d_const; i++ {
					// Bounds check
					row := y + j
					col := x + i
					if row >= 0 && row < 20 && col >= 0 && col < 40 {
						a := row*40 + col
						if a >= 0 && a < len(sim.wallBalls) {
							// current ball (y,x) is implicitly the center, but here we setAroundModule on 'a'?
							// Original:
							// int a = (y+j) * 40 + x+i;
							// Vec temp = new Vec(wallBalls.get(a)...);
							// wallBalls.get(a).setAroundMojule(temp);
							// Wait, setAroundModule adds to diff.
							// It gets its OWN position and adds to diff? That sounds weird.
							// `wallBalls.get(a).setAroundMojule(temp)`
							// `diff.add(v)`
							// It adds its own position to `diff`.
							// Then `moduleGravity`: `diff.mult(a); velocity.add(diff)`.
							// This looks like it's summing positions of neighbors?
							// But the loops iterate over the grid and for each cell, it iterates a window.
							// And for each neighbor 'a', it adds 'a's position to 'a's diff?
							// That means a cell is visited multiple times?
							// It accumulates position.

							target := sim.wallBalls[a]
							pos := NewVec(target.X, target.Y, target.Z)
							target.SetAroundModule(pos)
						}
					}
				}
			}
		}
	}

	for _, wb := range sim.wallBalls {
		wb.Gravity()
		wb.ModuleGravity()
		wb.Move()
	}
}

func main() {
	sim := NewSimulation()

	// Initial calibration
	fmt.Println("Calibrating...")
	sim.Calibration()
	fmt.Println("Calibration done.")

	fmt.Println("Starting simulation loop...")

	ticker := time.NewTicker(33 * time.Millisecond) // ~30 FPS
	defer ticker.Stop()

	// Run for a fixed amount of steps or infinite
	steps := 0
	for range ticker.C {
		sim.Recognize()
		sim.MakeForce()
		sim.Collision()
		sim.Move()
		sim.UpdateWallBalls()

		if steps%30 == 0 {
			fmt.Printf("Step %d: Forces count: %d\n", steps, len(sim.forces))
		}
		steps++

		if steps > 300 { // Run for ~10 seconds
			break
		}
	}
	fmt.Println("Simulation finished.")
}
