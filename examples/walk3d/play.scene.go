package main

import (
	"fmt"

	"github.com/MarcelArt/play-engine/engine"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type playScene struct {
	engine.BaseScene
	prefs  prefs
	world  *world
	player *player
}

func newPlayScene(p prefs) engine.Scene { return &playScene{prefs: p} }

func (s *playScene) Enter(ctx *engine.Context) {
	s.world = newWorld()
	s.player = newPlayer(s.prefs)
	ctx.SetMouseCaptured(true)
}

func (s *playScene) Exit(ctx *engine.Context) { ctx.SetMouseCaptured(false) }

func (s *playScene) Update(ctx *engine.Context, dt float32) {
	if ctx.Input.Pressed(actionPause) {
		ctx.Scenes.Push(scenePause)
		return
	}
	if ctx.Input.Pressed(actionFullscreen) {
		d := ctx.Display()
		d.Fullscreen = !d.Fullscreen
		ctx.ApplyDisplay(d)
		savePrefs(ctx, s.prefs)
	}
	s.player.Update(ctx, s.world, dt)
}

func (s *playScene) Draw(ctx *engine.Context) {
	s.player.cam.Begin()
	s.world.Draw()
	s.player.cam.End()

	// HUD, in screen space after the 3D pass.
	w, h := int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight())
	rl.DrawLine(w/2-8, h/2, w/2+8, h/2, rl.White)
	rl.DrawLine(w/2, h/2-8, w/2, h/2+8, rl.White)

	rl.DrawText("WASD move  Shift sprint  Space jump  Mouse look  F fullscreen  P pause  Esc quit", 10, h-26, 16, rl.White)
	f := s.player.feet
	rl.DrawText(fmt.Sprintf("x %.1f  y %.1f  z %.1f", f.X, f.Y, f.Z), 10, 24, 16, rl.White)
}
