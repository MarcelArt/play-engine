package main

import (
	"github.com/MarcelArt/play-engine/engine"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// drawBackground stretches the background over the whole window.
func drawBackground(ctx *engine.Context) {
	tex := ctx.Assets.Texture(texBackground)
	rl.DrawTexturePro(
		tex,
		rl.NewRectangle(0, 0, float32(tex.Width), float32(tex.Height)),
		rl.NewRectangle(0, 0, float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())),
		rl.Vector2Zero(),
		0,
		rl.White,
	)
}

// drawSprite draws tex centered on pos at size×size pixels, rotated by
// rotation radians. The sprites are square with the object centered, so the
// center is also the rotation point.
func drawSprite(tex rl.Texture2D, pos rl.Vector2, size, rotation float32) {
	rl.DrawTexturePro(
		tex,
		rl.NewRectangle(0, 0, float32(tex.Width), float32(tex.Height)),
		rl.NewRectangle(pos.X, pos.Y, size, size),
		rl.NewVector2(size/2, size/2),
		rotation*rl.Rad2deg,
		rl.White,
	)
}

func drawCentered(text string, y, size int32) {
	x := (int32(rl.GetScreenWidth()) - rl.MeasureText(text, size)) / 2
	rl.DrawText(text, x, y, size, rl.White)
}
