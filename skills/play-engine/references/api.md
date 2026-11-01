# play-engine API reference

Module: `github.com/MarcelArt/play-engine`

## engine

`engine.New(title string, opts ...engine.Option) *Game`

| Option | Effect |
|---|---|
| `WithSize(w, h int32)` | window size, default 1280×720 |
| `WithFPS(n int32)` | target FPS, default 60 |
| `WithResizable()` | resizable window |
| `WithFPSCounter()` | draw FPS on top |
| `WithBackground(rl.Color)` | clear color, default black |
| `WithExitKey(key int32)` | quit key, `0` disables, default Escape |
| `WithFixedStep(seconds float32)` | `FixedUpdate` interval, default 1/60 |
| `WithAssets(assets.Manifest)` | loaded after the window opens, before the first scene |
| `WithBindings(map[input.Action]input.Binding)` | input bindings |

- `(*Game).Register(name string, f Factory) *Game`; call before `Run`. Registering a name twice panics.
- `(*Game).Run(start string) error`; error only when an asset fails to load. Unknown scene names panic.
- `(*Game).Context() *Context`; services before `Run`.
- Fixed step: leftover time carries over, at most 5 catch-up steps per frame.
- Popping the last scene quits the game.

`Scene` methods: `Enter(ctx)`, `Update(ctx, dt)`, `FixedUpdate(ctx, dt)`, `Draw(ctx)`, `Exit(ctx)`. `Factory` is `func(ctx *Context) Scene`, run on every switch or push so each visit starts fresh.

`Context`: `Name`, `Assets *assets.Manager`, `Input *input.Map`, `Scenes *SceneManager`, `Quit()`.

`SceneManager`: `Switch(name)`, `Push(name)`, `Pop()`; all queued until the end of the frame.

## assets

```go
assets.Manifest{
    Textures: []assets.TextureSpec{{Name, Path, Data, Smooth}},
    Fonts:    []assets.FontSpec{{Name, Path}},
}
```

- `Data` (from `go:embed`) skips the disk read; `Path` still supplies the extension. Without `Data`, `Path` is relative to the working directory.
- `Smooth: true` → mipmaps + trilinear. Default is point filtering.
- `Manager`: `Texture(name)`, `Font(name)` (panic naming the asset when missing), `LoadTexture`, `LoadFont`, `UnloadAll`.
- Fonts load from a path only.

## input

- `input.Action` is a string type. `input.Binding{Keys []int32, MouseButtons []rl.MouseButton, GamepadButtons []int32}`; `input.Keys(keys...)` is the keyboard shorthand. Gamepad 0 is read.
- `Map`: `Down`, `Pressed`, `Released`, `Axis(neg, pos) float32` (−1, 0, 1), `Bind(action, binding)`.

## save

- `save.New[T](app, file) (*Store[T], error)` stores JSON under the user data dir (`~/.local/share/<app>/<file>` on Linux); `save.NewAt[T](path)` takes an explicit path.
- `Load() (T, error)` returns `save.ErrNotFound` when absent and a different error when corrupt. `LoadOr(def T) T`. `Save(v T) error` writes atomically. `Path()`.

## entity

- `Entity`: `Update(dt)`, `Draw()`, `Alive() bool`. Embed `entity.Lifecycle` for `Alive` and `Kill`.
- `World[T Entity]`: `Add`, `Update(dt)`, `Draw()`, `Sweep()`, `Items()`, `Len()`, `Clear()`. Entities added during `Update` first run on the next call; `Update` sweeps dead entities at its end. `Items()` is the live slice; do not keep it across `Update`.

## option

`option.Func[T]` and `option.Apply(&t, opts...)`; the pattern behind `engine.Option`. Use it for your own constructors.

## Examples

- `examples/asteroids`: wireframe, no assets.
- `examples/asteroids_2d`: sprites and background from embedded PNGs.
