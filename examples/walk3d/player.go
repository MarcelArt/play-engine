package main

import (
	"github.com/MarcelArt/play-engine/camera"
	"github.com/MarcelArt/play-engine/engine"
	"github.com/MarcelArt/play-engine/input"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	eyeHeight  float32 = 1.7
	radius     float32 = 0.4
	stepHeight float32 = 0.55 // blocks lower than this are walked up
	walkSpeed  float32 = 4.5  // m/s
	sprintMult float32 = 1.7
	jumpSpeed  float32 = 6
	gravity    float32 = 16
	stickTurn  float32 = 2.5 // rad/s at full right-stick deflection
)

// player is the first-person controller. It owns the camera and moves in
// Update, because it follows input directly.
type player struct {
	cam      *camera.FirstPerson
	feet     rl.Vector3
	velY     float32
	grounded bool

	// wasCaptured lets the first frame after capturing the mouse be ignored,
	// because the cursor jumps to the window center and reports a large delta.
	wasCaptured bool
}

func newPlayer(p prefs) *player {
	pl := &player{
		cam:  camera.NewFirstPerson(rl.NewVector3(0, eyeHeight, 0), p.FOV),
		feet: rl.NewVector3(0, 0, 0),
	}
	pl.cam.Sensitivity = p.Sensitivity
	return pl
}

func (p *player) Update(ctx *engine.Context, w *world, dt float32) {
	p.look(ctx, dt)
	p.move(ctx, w, dt)
	p.cam.SetEye(rl.Vector3Add(p.feet, rl.NewVector3(0, eyeHeight, 0)))
}

func (p *player) look(ctx *engine.Context, dt float32) {
	if ctx.MouseCaptured() {
		if p.wasCaptured {
			d := ctx.Input.MouseDelta()
			p.cam.Look(d.X, d.Y)
		}
		p.wasCaptured = true
	} else {
		p.wasCaptured = false
	}

	// The right stick turns at a fixed speed, scaled by dt.
	s := ctx.Input.Stick(input.RightStick)
	p.cam.Turn(-s.X*stickTurn*dt, -s.Y*stickTurn*dt)
}

func (p *player) move(ctx *engine.Context, w *world, dt float32) {
	in := ctx.Input

	// Keys and the left stick both feed forward/strafe in [-1, 1].
	stick := in.Stick(input.LeftStick)
	forward := in.Axis(actionBack, actionForward) - stick.Y
	strafe := in.Axis(actionLeft, actionRight) + stick.X

	wish := rl.Vector3Add(
		rl.Vector3Scale(p.cam.Forward(), forward),
		rl.Vector3Scale(p.cam.Right(), strafe),
	)
	if l := rl.Vector3Length(wish); l > 1 {
		wish = rl.Vector3Scale(wish, 1/l) // diagonals are not faster
	}
	speed := walkSpeed
	if in.Down(actionSprint) {
		speed *= sprintMult
	}
	p.feet.X += wish.X * speed * dt
	p.feet.Z += wish.Z * speed * dt
	w.pushOut(&p.feet, radius, stepHeight)

	// Vertical: gravity, jump, and landing on the highest surface underneath.
	if p.grounded && in.Pressed(actionJump) {
		p.velY = jumpSpeed
		p.grounded = false
	}
	prevY := p.feet.Y
	p.velY -= gravity * dt
	p.feet.Y += p.velY * dt

	// A surface counts as "underneath" if the feet were at most stepHeight below it.
	ground := w.groundAt(p.feet.X, p.feet.Z, max(prevY, p.feet.Y)+stepHeight)
	if p.feet.Y <= ground && p.velY <= 0 {
		p.feet.Y = ground
		p.velY = 0
		p.grounded = true
	} else if p.feet.Y > ground+0.01 {
		p.grounded = false
	}
}
