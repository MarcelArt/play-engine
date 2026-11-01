package main

import (
	"github.com/MarcelArt/play-engine/engine"
	"github.com/MarcelArt/play-engine/entity"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	shipSize            float32 = 20
	shipCollisionRadius         = shipSize * 0.7
	shipRotationSpeed   float32 = 4   // rad/s
	shipThrust          float32 = 350 // px/s²
	shipDrag            float32 = 0.5 // fraction of velocity lost per second
	shipMaxSpeed        float32 = 400 // px/s
)

// ship reads input from the engine's input map in Update and moves in
// FixedUpdate-sized steps; it reports a shot request instead of spawning
// bullets itself, so the scene stays in charge of the world.
type ship struct {
	entity.Lifecycle
	ctx      *engine.Context
	pos, vel rl.Vector2
	rotation float32

	wantsShot bool
}

func newShip(ctx *engine.Context, pos rl.Vector2) *ship {
	return &ship{ctx: ctx, pos: pos}
}

// Nose is the tip of the ship, where bullets spawn.
func (s *ship) Nose() rl.Vector2 { return s.local(rl.NewVector2(0, -shipSize)) }

// TakeShot reports whether a shot was requested since the last call.
func (s *ship) TakeShot() bool {
	shot := s.wantsShot
	s.wantsShot = false
	return shot
}

func (s *ship) Update(dt float32) {
	in := s.ctx.Input

	s.rotation += in.Axis(actionLeft, actionRight) * shipRotationSpeed * dt

	if in.Down(actionThrust) {
		dir := rl.Vector2Rotate(rl.NewVector2(0, -1), s.rotation)
		s.vel = rl.Vector2Add(s.vel, rl.Vector2Scale(dir, shipThrust*dt))
	}
	s.vel = rl.Vector2Scale(s.vel, 1-shipDrag*dt)
	if speed := rl.Vector2Length(s.vel); speed > shipMaxSpeed {
		s.vel = rl.Vector2Scale(s.vel, shipMaxSpeed/speed)
	}
	s.pos = rl.Vector2Add(s.pos, rl.Vector2Scale(s.vel, dt))

	if in.Pressed(actionShoot) {
		s.wantsShot = true
	}
}

func (s *ship) Draw() {
	rl.DrawTriangleLines(
		s.local(rl.NewVector2(0, -shipSize)),
		s.local(rl.NewVector2(-shipSize*0.6, shipSize)),
		s.local(rl.NewVector2(shipSize*0.6, shipSize)),
		rl.White,
	)
}

func (s *ship) local(v rl.Vector2) rl.Vector2 {
	return rl.Vector2Add(s.pos, rl.Vector2Rotate(v, s.rotation))
}
