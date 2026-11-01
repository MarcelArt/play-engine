package main

import (
	"fmt"

	"github.com/MarcelArt/play-engine/engine"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type menuScene struct {
	engine.BaseScene
	highScore uint
}

func newMenuScene(ctx *engine.Context) engine.Scene {
	s := &menuScene{}
	if store := highScoreStore(ctx); store != nil {
		s.highScore = store.LoadOr(saveFile{}).HighScore
	}
	return s
}

func (s *menuScene) Update(ctx *engine.Context, dt float32) {
	if ctx.Input.Pressed(actionSelect) {
		ctx.Scenes.Switch(scenePlay)
	}
}

func (s *menuScene) Draw(ctx *engine.Context) {
	h := rl.GetScreenHeight()
	drawCentered("ASTEROIDS", int32(h)/2-80, 48)
	drawCentered("Enter to play", int32(h)/2, 24)
	drawCentered(fmt.Sprintf("High score: %d", s.highScore), int32(h)/2+40, 20)
	drawCentered("Arrows/WASD move, Space shoot, P pause, Esc quit", int32(h)-40, 16)
}

func drawCentered(text string, y, size int32) {
	x := (int32(rl.GetScreenWidth()) - rl.MeasureText(text, size)) / 2
	rl.DrawText(text, x, y, size, rl.White)
}
