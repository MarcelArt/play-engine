package main

import (
	"fmt"

	"github.com/MarcelArt/play-engine/engine"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type menuScene struct {
	engine.BaseScene
	highScore uint
	time      float32
}

func newMenuScene(ctx *engine.Context) engine.Scene {
	s := &menuScene{}
	if store := highScoreStore(ctx); store != nil {
		s.highScore = store.LoadOr(saveFile{}).HighScore
	}
	return s
}

func (s *menuScene) Update(ctx *engine.Context, dt float32) {
	s.time += dt
	if ctx.Input.Pressed(actionSelect) {
		ctx.Scenes.Switch(scenePlay)
	}

	// Each change applies to the open window and is saved to the config file.
	if ctx.Input.Pressed(actionFullscreen) {
		changeDisplay(ctx, func(d *engine.Display) { d.Fullscreen = !d.Fullscreen })
	}
	if ctx.Input.Pressed(actionVSync) {
		changeDisplay(ctx, func(d *engine.Display) { d.VSync = !d.VSync })
	}
	if ctx.Input.Pressed(actionResolution) {
		changeDisplay(ctx, func(d *engine.Display) {
			r := engine.NextResolution(d.Width, d.Height)
			d.Width, d.Height = r.Width, r.Height
		})
	}
}

func (s *menuScene) Draw(ctx *engine.Context) {
	drawBackground(ctx)

	w, h := float32(rl.GetScreenWidth()), int32(rl.GetScreenHeight())
	// A slowly spinning ship as the title art.
	drawSprite(ctx.Assets.Texture(texShip), rl.NewVector2(w/2, float32(h)/2-110), 120, s.time*0.5)

	drawCentered("ASTEROIDS", h/2-30, 48)
	drawCentered("Enter to play", h/2+30, 24)
	drawCentered(fmt.Sprintf("High score: %d", s.highScore), h/2+70, 20)
	drawCentered("Arrows/WASD move, Space shoot, P pause, Esc quit", h-60, 16)

	d := ctx.Display()
	drawCentered(fmt.Sprintf("F fullscreen: %v   V vsync: %v   R resolution: %dx%d",
		d.Fullscreen, d.VSync, d.Width, d.Height), h-34, 16)
}
