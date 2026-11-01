package main

import (
	"math"

	"github.com/MarcelArt/play-engine/entity"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	rockRadius   float32 = 40
	rockVertices         = 10
	rockSpeed    float32 = 80
)

type rock struct {
	entity.Lifecycle
	pos, vel rl.Vector2
	// shape is a per-vertex radius multiplier in [0.7, 1.0], rolled once so
	// every rock keeps its own jagged outline.
	shape [rockVertices]float32
}

func newRock(pos, target rl.Vector2) *rock {
	r := &rock{pos: pos}
	dir := rl.Vector2Normalize(rl.Vector2Subtract(target, pos))
	r.vel = rl.Vector2Scale(dir, rockSpeed)
	for i := range r.shape {
		r.shape[i] = 0.7 + float32(rl.GetRandomValue(0, 30))/100
	}
	return r
}

func (r *rock) Update(dt float32) {
	r.pos = rl.Vector2Add(r.pos, rl.Vector2Scale(r.vel, dt))
}

func (r *rock) Draw() {
	var pts [rockVertices + 1]rl.Vector2
	for i := range pts {
		k := i % rockVertices
		a := 2 * math.Pi * float64(k) / rockVertices
		d := rockRadius * r.shape[k]
		pts[i] = rl.NewVector2(r.pos.X+float32(math.Sin(a))*d, r.pos.Y+float32(math.Cos(a))*d)
	}
	rl.DrawLineStrip(pts[:], rl.White)
}
