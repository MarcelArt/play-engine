---
name: play-engine-new-game
description: Scaffold a new Go raylib game on play-engine (github.com/MarcelArt/play-engine). Use when the user starts a new game, prototype or project with play-engine, or asks for a play-engine project layout.
---

# New play-engine game

Load the `play-engine` skill for the engine's concepts and API; this skill only creates the project.

## Steps

1. Ask for the game name and Go module path if they are not given. The module path is the one decision to ask for.
2. Create the folder and initialise the module:
   ```bash
   mkdir <name> && cd <name> && go mod init <module>
   go get github.com/MarcelArt/play-engine@latest
   ```
   When `play-engine` is checked out beside the game and is unpublished, use the checkout instead:
   `go mod edit -replace github.com/MarcelArt/play-engine=../play-engine`, then `go mod tidy`.
3. Lay out the files:
   ```
   main.go
   internal/scenes/menu.scene.go
   internal/scenes/play.scene.go
   internal/entities/player.entity.go
   assets/
   ```
   Scenes and entities live in `internal/`; `main.go` only wires. For a single-scene prototype, `main.go` plus one `play.scene.go` in package `main` is enough.
4. Write `main.go`:
   ```go
   package main

   import (
   	"log"

   	"github.com/MarcelArt/play-engine/engine"
   	"github.com/MarcelArt/play-engine/input"
   	rl "github.com/gen2brain/raylib-go/raylib"
   )

   const (
   	sceneMenu = "menu"
   	scenePlay = "play"
   )

   const (
   	actionLeft  input.Action = "left"
   	actionRight input.Action = "right"
   )

   func main() {
   	g := engine.New("<name>",
   		engine.WithSize(1280, 720),
   		engine.WithFPSCounter(),
   		engine.WithBindings(map[input.Action]input.Binding{
   			actionLeft:  input.Keys(rl.KeyLeft, rl.KeyA),
   			actionRight: input.Keys(rl.KeyRight, rl.KeyD),
   		}),
   	)
   	g.Register(sceneMenu, newMenuScene)
   	g.Register(scenePlay, newPlayScene)

   	if err := g.Run(sceneMenu); err != nil {
   		log.Fatal(err)
   	}
   }
   ```
   When scenes sit in `internal/scenes`, export the name constants and factories there and import them in `main.go`.
5. Write a menu scene that calls `ctx.Scenes.Switch(scenePlay)` on a key press, and a play scene that draws something and moves it with `ctx.Input.Axis(actionLeft, actionRight)`.
6. Add assets with the `play-engine` skill's "Add a texture" steps, embedding files that sit inside the package directory that embeds them.
7. Run `gofmt -l .`, `go vet ./...`, `go build ./...`, then `go run .`.

Done when `go run .` opens the window on the menu and the key press reaches the play scene, and vet and build pass. Running needs a graphical environment; when none is available, report that the run was not verified.
