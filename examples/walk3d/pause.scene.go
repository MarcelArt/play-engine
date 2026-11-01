package main

import (
	"github.com/MarcelArt/play-engine/engine"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// pauseScene is pushed over the play scene. It frees the mouse while it is
// open and captures it again when it closes; the play scene below keeps
// drawing but stops updating.
type pauseScene struct{ engine.BaseScene }

func newPauseScene(*engine.Context) engine.Scene { return &pauseScene{} }

func (s *pauseScene) Enter(ctx *engine.Context) { ctx.SetMouseCaptured(false) }
func (s *pauseScene) Exit(ctx *engine.Context)  { ctx.SetMouseCaptured(true) }

func (s *pauseScene) Update(ctx *engine.Context, dt float32) {
	if ctx.Input.Pressed(actionPause) {
		ctx.Scenes.Pop()
	}
}

func (s *pauseScene) Draw(ctx *engine.Context) {
	w, h := int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight())
	rl.DrawRectangle(0, 0, w, h, rl.NewColor(0, 0, 0, 140))
	text := "PAUSED - P to resume"
	rl.DrawText(text, (w-rl.MeasureText(text, 32))/2, h/2-16, 32, rl.White)
}
