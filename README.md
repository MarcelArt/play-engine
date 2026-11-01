# play-engine

A small scene-based game loop on top of [raylib-go](https://github.com/gen2brain/raylib-go). It handles the window, a fixed-timestep loop, a scene stack, assets, input actions and saves, so each game only contains its own rules.

You keep using `rl.` directly for drawing. The engine does not wrap raylib.

```
go get github.com/MarcelArt/play-engine
```

Requires Go 1.27 and raylib's usual native dependencies (a graphical environment to run).

## Quick start

```go
package main

import (
	"log"

	"github.com/MarcelArt/play-engine/engine"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type hello struct{ engine.BaseScene }

func (hello) Draw(ctx *engine.Context) {
	rl.DrawText("hello", 20, 20, 32, rl.White)
}

func main() {
	g := engine.New("hello", engine.WithSize(800, 540))
	g.Register("hello", func(ctx *engine.Context) engine.Scene { return hello{} })
	if err := g.Run("hello"); err != nil {
		log.Fatal(err)
	}
}
```

Two full games live in `examples/`:

- [`examples/asteroids`](examples/asteroids): wireframe shapes, no asset files. Run it with `go run ./examples/asteroids`.
- [`examples/asteroids_2d`](examples/asteroids_2d): the same game with sprites and a background, loaded from embedded PNG files. Run it with `go run ./examples/asteroids_2d`.

## Agent skills

The repo ships two skills in [`skills/`](skills) that teach coding agents (Claude Code, Cursor, Codex and others) how to use the engine. Install them with [`npx skills`](https://github.com/vercel-labs/skills):

```bash
npx skills add MarcelArt/play-engine --list          # show what is available
npx skills add MarcelArt/play-engine                 # install interactively
npx skills add MarcelArt/play-engine --skill play-engine -a claude-code -g -y
```

| Skill | Use |
|---|---|
| `play-engine` | Add a scene, entity, texture, input action or save to an existing play-engine game, with the API reference |
| `play-engine-new-game` | Scaffold a new game project on play-engine |

`-g` installs for your user instead of the current project, and `-a` picks the agent. To try them from a local checkout, run `npx skills add ./play-engine` from the folder above it. The repo must be pushed to GitHub for the `MarcelArt/play-engine` form to work.

## Concepts

### Game and options

`engine.New(title, opts...)` configures a game. Nothing touches raylib until `Run`.

| Option | Effect |
|---|---|
| `WithSize(w, h)` | Window size (default 1280×720) |
| `WithFPS(n)` | Target FPS (default 60) |
| `WithResizable()` | Allow resizing the window |
| `WithFPSCounter()` | Draw the FPS counter on top |
| `WithBackground(color)` | Clear color (default black) |
| `WithExitKey(key)` | Key that quits; `0` disables it (default Escape) |
| `WithFixedStep(seconds)` | `FixedUpdate` interval (default 1/60) |
| `WithAssets(manifest)` | Assets loaded before the first scene |
| `WithBindings(map)` | Input action bindings |

`Run(start)` opens the window, loads assets, enters the start scene and loops until the game quits. It returns an error if an asset fails to load. Cleanup runs even on panic.

### Scenes

A scene is one screen: menu, level, pause overlay. Implement `engine.Scene`, or embed `engine.BaseScene` and override only what you need.

| Method | When it runs |
|---|---|
| `Enter(ctx)` | Once, after the window exists, when the scene joins the stack. Load things and build state here, not in a constructor. |
| `Update(ctx, dt)` | Every frame, top scene only. Read input and move things here. |
| `FixedUpdate(ctx, dt)` | Zero or more times per frame at the fixed step, top scene only. Put simulation and physics here. |
| `Draw(ctx)` | Every frame, for every scene in the stack, bottom first. |
| `Exit(ctx)` | Once, when the scene leaves the stack. |

Register a factory for each scene by name before `Run`. The factory runs on every visit, so each visit starts with fresh state.

```go
g.Register("menu", newMenuScene)
g.Register("play", newPlayScene)
```

Move between scenes through the context. Scenes never import each other, only the names:

```go
ctx.Scenes.Switch("play") // replace the whole stack
ctx.Scenes.Push("pause")  // put a scene on top; the one below keeps drawing but stops updating
ctx.Scenes.Pop()          // remove the top scene; popping the last one quits
ctx.Quit()                // end the game after this frame
```

These calls are queued and applied between frames, so it is safe to call them from inside `Update`.

### Context

Every scene method receives a `*engine.Context`, which holds the game's services:

- `ctx.Name`: the game title, used to namespace saves
- `ctx.Assets`: loaded textures and fonts
- `ctx.Input`: action-based input
- `ctx.Scenes`: the scene manager
- `ctx.Quit()`

Pass entities only what they need (a texture, the input map) in their constructor, not the whole context. That keeps them easy to test.

### Assets

List what to load in a manifest, then fetch by name in a scene's `Enter`:

```go
engine.New("game", engine.WithAssets(assets.Manifest{
	Textures: []assets.TextureSpec{
		{Name: "player", Path: "assets/player.png"},
		{Name: "photo", Path: "assets/photo.png", Smooth: true},
	},
	Fonts: []assets.FontSpec{{Name: "ui", Path: "assets/ui.ttf"}},
}))

tex := ctx.Assets.Texture("player")
```

- Textures use point filtering by default, which keeps pixel art crisp. Set `Smooth: true` for mipmapped trilinear filtering, which looks cleaner when large art is drawn small.
- A missing file makes `Run` return an error naming the asset.
- Asking for a name that was never loaded panics with that name, so a typo fails loudly instead of drawing nothing.
- Assets are freed when the game exits.

#### Embedded assets

Paths like `assets/player.png` are read from disk relative to the folder you run the game from, so the game breaks if you start it somewhere else. To avoid that, embed the files into the program with Go's `go:embed` and pass their bytes in `Data`:

```go
import _ "embed"

//go:embed assets/player.png
var playerPNG []byte

engine.WithAssets(assets.Manifest{
	Textures: []assets.TextureSpec{
		{Name: "player", Path: "player.png", Data: playerPNG, Smooth: true},
	},
})
```

- When `Data` is set, the file is not read from disk. `Path` is then only used for its extension (`.png`), which tells raylib the image format, so keep it.
- The files are compiled into the binary, so you can ship one executable with no `assets/` folder next to it.
- The `go:embed` line must sit directly above the variable, and the file must be inside the game's own package directory or a subfolder of it (embed cannot reach `../`).
- Fonts are not embeddable yet. `FontSpec` still loads from a path.

### Input actions

Game code asks about actions and never names a key. Define the actions once and bind them in `main`:

```go
const (
	actionJump  input.Action = "jump"
	actionLeft  input.Action = "left"
	actionRight input.Action = "right"
)

engine.WithBindings(map[input.Action]input.Binding{
	actionJump:  input.Keys(rl.KeySpace),
	actionLeft:  input.Keys(rl.KeyLeft, rl.KeyA),
	actionRight: input.Keys(rl.KeyRight, rl.KeyD),
})
```

A `Binding` can also list `MouseButtons` and `GamepadButtons` (gamepad 0). Use it in a scene:

```go
ctx.Input.Down(actionJump)                  // held
ctx.Input.Pressed(actionJump)               // went down this frame
ctx.Input.Released(actionJump)              // went up this frame
ctx.Input.Axis(actionLeft, actionRight)     // -1, 0 or 1
ctx.Input.Bind(actionJump, input.Keys(rl.KeyW)) // rebind at runtime
```

### Saves

`save.Store[T]` stores one value as JSON under the user data directory (`~/.local/share/<app>/<file>` on Linux), namespaced by game so games never overwrite each other.

```go
type saveFile struct{ HighScore uint }

store, err := save.New[saveFile](ctx.Name, "save.json")
data := store.LoadOr(saveFile{})   // default when missing or unreadable
err = store.Save(saveFile{HighScore: 42})
```

`Load()` returns `save.ErrNotFound` when nothing has been saved yet, and a different error for a corrupt file. Writes are atomic: a crash never leaves a half-written file. For tests, `save.NewAt[T](path)` uses an explicit path.

### Entities (optional)

`entity.World[T]` keeps a list of things to update and draw. Use it when a scene has many short-lived objects such as bullets or enemies. Games that manage their own data, like a tile grid, can skip it.

An entity implements `Update(dt)`, `Draw()` and `Alive() bool`. Embed `entity.Lifecycle` to get `Alive` and `Kill`:

```go
type bullet struct {
	entity.Lifecycle
	pos rl.Vector2
}

func (b *bullet) Update(dt float32) { /* move; call b.Kill() when done */ }
func (b *bullet) Draw()             { /* draw */ }

var bullets entity.World[*bullet]
bullets.Add(&bullet{})
bullets.Update(dt) // updates live entities, then removes dead ones
bullets.Draw()
```

- Entities added during `Update` are first updated on the next call.
- Killing entities during collision handling is safe. Call `Sweep()` afterwards to remove them.

### Options helper

`option.Func[T]` and `option.Apply` are the functional-options helper the library uses. Use them for your own constructors so every type in your game is configured the same way:

```go
type Config struct{ Speed float32 }
type Option = option.Func[Config]

func WithSpeed(s float32) Option { return func(c *Config) { c.Speed = s } }

cfg := option.Apply(&Config{}, WithSpeed(3))
```

## Fixed vs variable update

`Update` runs once per rendered frame with the real frame time. `FixedUpdate` runs at a constant step, so simulation behaves the same at any frame rate. Leftover time carries into the next frame, and at most five catch-up steps run per frame so a long stall cannot spiral.

Rule of thumb: read input in `Update`, run physics and simulation in `FixedUpdate`.

## Suggested game layout

```
mygame/
  main.go                 // engine.New, bindings, Register, Run
  internal/
    scenes/
      menu.scene.go
      play.scene.go
    entities/
      player.entity.go
```

`main.go` is the one place that lists every scene and binding, so reading it shows the whole shape of the game.

## Using it from a local checkout

While developing the engine and a game together, point the game at your local copy:

```
go mod edit -replace github.com/MarcelArt/play-engine=../play-engine
```

or use a `go.work` file.

## Development

```
go vet ./...
go test ./...
```
