// Asteroids is the example game for play-engine: menu, play and pause scenes,
// action-based input, a fixed-step simulation, an entity world and a saved
// high score.
package main

import (
	"log"

	"github.com/MarcelArt/play-engine/engine"
	"github.com/MarcelArt/play-engine/input"
	rl "github.com/gen2brain/raylib-go/raylib"
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
	g := engine.New("asteroids",
		engine.WithSize(800, 540),
		engine.WithFPSCounter(),
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
