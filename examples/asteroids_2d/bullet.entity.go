package main

import (
	"github.com/MarcelArt/play-engine/entity"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	bulletSpeed  float32 = 300
	bulletRadius float32 = 3
	bulletLife   float32 = 2
)

type bullet struct {
	entity.Lifecycle
	pos, vel rl.Vector2
	age      float32
}

func newBullet(pos rl.Vector2, rotation float32) *bullet {
	dir := rl.Vector2Rotate(rl.NewVector2(0, -1), rotation)
	return &bullet{pos: pos, vel: rl.Vector2Scale(dir, bulletSpeed)}
}

func (b *bullet) Update(dt float32) {
	b.pos = rl.Vector2Add(b.pos, rl.Vector2Scale(b.vel, dt))
	b.age += dt
	if b.age > bulletLife {
		b.Kill()
	}
}

func (b *bullet) Draw() { rl.DrawCircleV(b.pos, bulletRadius, rl.White) }
