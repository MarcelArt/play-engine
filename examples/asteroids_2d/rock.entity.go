package main

import (
	"github.com/MarcelArt/play-engine/entity"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	rockSize   float32 = 110 // sprite size in pixels, including its transparent padding
	rockRadius float32 = 38  // collision radius; the visible rock is about 80% of the sprite
	rockSpeed  float32 = 80
)

type rock struct {
	entity.Lifecycle
	tex      rl.Texture2D
	pos, vel rl.Vector2
	rotation float32
	spin     float32 // radians per second
}

func newRock(tex rl.Texture2D, pos, target rl.Vector2) *rock {
	dir := rl.Vector2Normalize(rl.Vector2Subtract(target, pos))
	return &rock{
		tex:      tex,
		pos:      pos,
		vel:      rl.Vector2Scale(dir, rockSpeed),
		rotation: float32(rl.GetRandomValue(0, 628)) / 100,
		spin:     float32(rl.GetRandomValue(-100, 100)) / 100,
	}
}

func (r *rock) Update(dt float32) {
	r.pos = rl.Vector2Add(r.pos, rl.Vector2Scale(r.vel, dt))
	r.rotation += r.spin * dt
}

func (r *rock) Draw() { drawSprite(r.tex, r.pos, rockSize, r.rotation) }
