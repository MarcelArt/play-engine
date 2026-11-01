package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/MarcelArt/play-engine/engine"
	"github.com/MarcelArt/play-engine/entity"
	"github.com/MarcelArt/play-engine/save"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const rockSpawnInterval float32 = 2

type playScene struct {
	engine.BaseScene

	ship    *ship
	bullets entity.World[*bullet]
	rocks   entity.World[*rock]

	store     *save.Store[saveFile]
	score     uint
	highScore uint
	dead      bool
	spawnIn   float32
}

func newPlayScene(ctx *engine.Context) engine.Scene {
	s := &playScene{store: highScoreStore(ctx), spawnIn: rockSpawnInterval}
	if s.store != nil {
		s.highScore = s.store.LoadOr(saveFile{}).HighScore
	}
	return s
}

func (s *playScene) Enter(ctx *engine.Context) {
	center := rl.NewVector2(float32(rl.GetScreenWidth())/2, float32(rl.GetScreenHeight())/2)
	s.ship = newShip(ctx, center)
}

func (s *playScene) Update(ctx *engine.Context, dt float32) {
	if s.dead {
		if ctx.Input.Pressed(actionSelect) {
			ctx.Scenes.Switch(scenePlay)
		}
		return
	}
	if ctx.Input.Pressed(actionPause) {
		ctx.Scenes.Push(scenePause)
		return
	}
	s.ship.Update(dt)
	if s.ship.TakeShot() {
		s.bullets.Add(newBullet(s.ship.Nose(), s.ship.rotation))
	}
}

func (s *playScene) FixedUpdate(ctx *engine.Context, dt float32) {
	if s.dead {
		return
	}
	s.bullets.Update(dt)
	s.rocks.Update(dt)

	s.spawnIn -= dt
	if s.spawnIn <= 0 {
		s.spawnRock()
		s.spawnIn = rockSpawnInterval
	}

	s.collide()
}

func (s *playScene) Draw(ctx *engine.Context) {
	if s.dead {
		drawCentered("GAME OVER", int32(rl.GetScreenHeight())/2-50, 36)
		drawCentered(fmt.Sprintf("Score %d  High %d", s.score, s.highScore), int32(rl.GetScreenHeight())/2, 24)
		drawCentered("Enter to retry", int32(rl.GetScreenHeight())/2+40, 20)
		return
	}
	s.ship.Draw()
	s.bullets.Draw()
	s.rocks.Draw()
	rl.DrawText(fmt.Sprint(s.score), int32(rl.GetScreenWidth())/2, 0, 36, rl.White)
}

func (s *playScene) spawnRock() {
	const margin = 10
	w, h := float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())
	corners := []rl.Vector2{
		{X: -margin, Y: -margin}, {X: -margin, Y: h + margin},
		{X: w + margin, Y: -margin}, {X: w + margin, Y: h + margin},
	}
	s.rocks.Add(newRock(corners[rand.IntN(len(corners))], s.ship.pos))
}

func (s *playScene) collide() {
	for _, b := range s.bullets.Items() {
		for _, r := range s.rocks.Items() {
			if b.Alive() && r.Alive() && rl.CheckCollisionCircles(b.pos, bulletRadius, r.pos, rockRadius) {
				b.Kill()
				r.Kill()
				s.score++
			}
		}
	}
	s.bullets.Sweep()
	s.rocks.Sweep()

	for _, r := range s.rocks.Items() {
		if rl.CheckCollisionCircles(r.pos, rockRadius, s.ship.pos, shipCollisionRadius) {
			s.lose()
			return
		}
	}
}

func (s *playScene) lose() {
	s.dead = true
	if s.score > s.highScore {
		s.highScore = s.score
		if s.store != nil {
			if err := s.store.Save(saveFile{HighScore: s.score}); err != nil {
				fmt.Println("could not save high score:", err)
			}
		}
	}
}
