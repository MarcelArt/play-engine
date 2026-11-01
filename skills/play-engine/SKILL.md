---
name: play-engine
description: Use play-engine (github.com/MarcelArt/play-engine), the Go + raylib scene engine, when a project imports it or the user wants scenes, entities, assets, input actions or saves in a Go raylib game. Covers adding a scene, entity, texture or input action to an existing play-engine game.
---

# play-engine

Scene-based loop over raylib-go. Draw with `rl.` directly; the engine does not wrap raylib. Full API: [references/api.md](references/api.md).

## Model

- `engine.New(title, opts...)` → `Register(name, factory)` → `Run(start)`. `main.go` is the single place that lists scenes, bindings and assets.
- A scene is `Enter`, `Update`, `FixedUpdate`, `Draw`, `Exit`. Embed `engine.BaseScene`, override what you use.
- Every scene method receives `*engine.Context`: `Assets`, `Input`, `Scenes`, `Name`, `Quit()`. There are no package globals.
- `Update` reads input and moves things. `FixedUpdate` runs the simulation at a fixed step. Only the top scene gets either; every scene in the stack gets `Draw`, bottom first.
- `ctx.Scenes.Switch(name)` replaces the stack, `Push` overlays, `Pop` returns. Changes land between frames, so call them from inside `Update`.
- Scenes reference each other by name constant only, never by import.

## Add a scene

1. Create `<name>.scene.go` with a struct embedding `engine.BaseScene` and a `new<Name>Scene(ctx *engine.Context) engine.Scene` factory.
2. Build state in `Enter`; the window exists by then.
3. Add a name constant and `g.Register(name, factory)` in `main.go`.
4. Navigate to it with `ctx.Scenes.Switch/Push` from another scene.

Done when `Run` reaches the scene and `Exit` releases whatever `Enter` created.

A pushed overlay (pause, dialog) draws over the scene below, so it draws only its own layer and never clears the background.

## Add an entity

1. Create `<name>.entity.go` with `Update(dt)`, `Draw()` and an embedded `entity.Lifecycle` (gives `Alive`, `Kill`).
2. Hold entities in `entity.World[*T]` on the scene; call `world.Update(dt)` in `FixedUpdate` and `world.Draw()` in `Draw`.
3. Kill during collision handling, then call `world.Sweep()`.
4. Hand each entity only what it needs (a texture, the input map) through its constructor. An entity requests actions (a `TakeShot()` flag) and the scene performs them.

## Add a texture

1. Put the file in the game package's `assets/` folder.
2. Embed it: `//go:embed assets/x.png` above `var xPNG []byte`.
3. Add `{Name: texX, Path: "x.png", Data: xPNG, Smooth: true}` to the `assets.Manifest` in `main.go`. `Smooth` mipmaps large art drawn small; omit it for pixel art.
4. Fetch with `ctx.Assets.Texture(texX)` in `Enter` and keep it in the struct. A name that was never loaded panics with that name.

## Add an input action

1. Declare `const actionX input.Action = "x"`.
2. Bind it in `engine.WithBindings`: `actionX: input.Keys(rl.KeyX, rl.KeyY)`.
3. Query it with `ctx.Input.Down/Pressed/Released(actionX)` or `Axis(neg, pos)`; game code never names a key.

## Save data

`save.New[T](ctx.Name, "save.json")` returns a store namespaced by game. Use `LoadOr(default)` for first runs and check the error from `Save`.

## Conventions

- File suffixes: `*.scene.go`, `*.entity.go`; raylib imported as `rl`.
- Name constants for scenes, actions and textures; strings appear once.
- Run `gofmt`, `go vet ./...` and `go test ./...` after changes.
