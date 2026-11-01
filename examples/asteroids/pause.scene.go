package main

import (
	"github.com/MarcelArt/play-engine/engine"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// pauseScene is pushed on top of the play scene. The play scene stops
// updating but keeps drawing underneath, so the dim overlay shows the frozen game.
type pauseScene struct{ engine.BaseScene }

func newPauseScene(*engine.Context) engine.Scene { return &pauseScene{} }

func (s *pauseScene) Update(ctx *engine.Context, dt float32) {
	if ctx.Input.Pressed(actionPause) {
		ctx.Scenes.Pop()
	}
	if ctx.Input.Pressed(actionSelect) {
		ctx.Scenes.Switch(sceneMenu)
	}
}

func (s *pauseScene) Draw(ctx *engine.Context) {
	rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), rl.NewColor(0, 0, 0, 160))
	h := int32(rl.GetScreenHeight())
	drawCentered("PAUSED", h/2-40, 36)
	drawCentered("P to resume, Enter for menu", h/2+10, 20)
}
