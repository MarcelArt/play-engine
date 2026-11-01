// Asteroids 2D is the sprite version of the asteroids example. It does the
// same things, but draws embedded PNG assets through the engine's asset
// manager instead of wireframe shapes.
package main

import (
	_ "embed"
	"log"

	"github.com/MarcelArt/play-engine/assets"
	"github.com/MarcelArt/play-engine/engine"
	"github.com/MarcelArt/play-engine/input"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Embedding the files means the game runs from any working directory.
var (
	//go:embed assets/ship.png
	shipPNG []byte
	//go:embed assets/asteroid.png
	asteroidPNG []byte
	//go:embed assets/space-background.png
	backgroundPNG []byte
)

// Asset names, used with ctx.Assets.Texture(...).
const (
	texShip       = "ship"
	texAsteroid   = "asteroid"
	texBackground = "background"
)

// Scene names. Scenes refer to each other only through these.
const (
	sceneMenu  = "menu"
	scenePlay  = "play"
	scenePause = "pause"
)

// Input actions.
const (
	actionLeft   input.Action = "left"
	actionRight  input.Action = "right"
	actionThrust input.Action = "thrust"
	actionShoot  input.Action = "shoot"
	actionPause  input.Action = "pause"
	actionSelect input.Action = "select"
)

func main() {
	g := engine.New("asteroids_2d",
		engine.WithSize(800, 540),
		engine.WithFPSCounter(),
		engine.WithAssets(assets.Manifest{
			Textures: []assets.TextureSpec{
				{Name: texShip, Path: "ship.png", Data: shipPNG, Smooth: true},
				{Name: texAsteroid, Path: "asteroid.png", Data: asteroidPNG, Smooth: true},
				{Name: texBackground, Path: "space-background.png", Data: backgroundPNG, Smooth: true},
			},
		}),
		engine.WithBindings(map[input.Action]input.Binding{
			actionLeft:   input.Keys(rl.KeyLeft, rl.KeyA),
			actionRight:  input.Keys(rl.KeyRight, rl.KeyD),
			actionThrust: input.Keys(rl.KeyUp, rl.KeyW),
			actionShoot:  input.Keys(rl.KeySpace),
			actionPause:  input.Keys(rl.KeyP),
			actionSelect: input.Keys(rl.KeyEnter),
		}),
	)

	g.Register(sceneMenu, newMenuScene)
	g.Register(scenePlay, newPlayScene)
	g.Register(scenePause, newPauseScene)

	if err := g.Run(sceneMenu); err != nil {
		log.Fatal(err)
	}
}
