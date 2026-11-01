// Walk3D is the first-person example for play-engine: mouse look, WASD
// movement, jumping, a pause scene that releases the mouse, and persisted
// display settings.
package main

import (
	"log"

	"github.com/MarcelArt/play-engine/engine"
	"github.com/MarcelArt/play-engine/input"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Scene names.
const (
	scenePlay  = "play"
	scenePause = "pause"
)

// Input actions.
const (
	actionForward    input.Action = "forward"
	actionBack       input.Action = "back"
	actionLeft       input.Action = "left"
	actionRight      input.Action = "right"
	actionJump       input.Action = "jump"
	actionSprint     input.Action = "sprint"
	actionPause      input.Action = "pause"
	actionFullscreen input.Action = "fullscreen"
)

func main() {
	p := loadPrefs()

	g := engine.New(appName,
		engine.WithDisplay(p.Display),
		engine.WithMSAA(),
		engine.WithFPSCounter(),
		engine.WithBackground(rl.NewColor(135, 190, 235, 255)), // sky
		engine.WithBindings(map[input.Action]input.Binding{
			actionForward:    input.Keys(rl.KeyW, rl.KeyUp),
			actionBack:       input.Keys(rl.KeyS, rl.KeyDown),
			actionLeft:       input.Keys(rl.KeyA, rl.KeyLeft),
			actionRight:      input.Keys(rl.KeyD, rl.KeyRight),
			actionJump:       {Keys: []int32{rl.KeySpace}, GamepadButtons: []int32{rl.GamepadButtonRightFaceDown}},
			actionSprint:     input.Keys(rl.KeyLeftShift),
			actionPause:      input.Keys(rl.KeyP),
			actionFullscreen: input.Keys(rl.KeyF),
		}),
	)

	g.Register(scenePlay, func(ctx *engine.Context) engine.Scene { return newPlayScene(p) })
	g.Register(scenePause, newPauseScene)

	if err := g.Run(scenePlay); err != nil {
		log.Fatal(err)
	}
}
