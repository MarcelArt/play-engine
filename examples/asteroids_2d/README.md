# asteroids_2d

The same game as [`asteroids`](../asteroids), drawn with sprites and a background image instead of wireframes. Read it to see how assets are embedded, loaded and drawn, and how display settings are changed and saved.

```bash
go run ./examples/asteroids_2d      # from the repository root
```

| Key | Action |
|---|---|
| Enter | start / retry |
| Arrows or A/D | rotate |
| Up or W | thrust |
| Space | shoot |
| P | pause / resume |
| F, V, R (on the menu) | toggle fullscreen, toggle vsync, cycle resolution |
| Esc | quit |

## What it shows

- Embedding image files with `go:embed` so the game runs from any folder
- Loading assets through a manifest and fetching them by name
- Drawing rotated sprites
- Changing the window at runtime with `ctx.ApplyDisplay` and saving it with the `settings` package

For the scene, entity and collision logic, which is unchanged, read the [`asteroids` README](../asteroids/README.md) first. This one covers only what differs.

## Files

| File | Role |
|---|---|
| `main.go` | Embeds the images, builds the manifest, binds actions, registers scenes |
| `draw.go` | `drawBackground`, `drawSprite` and `drawCentered` helpers |
| `settings.go` | The saved settings type and `changeDisplay` |
| `menu.scene.go` | Title screen with a spinning ship and the display toggles |
| `play.scene.go`, `pause.scene.go` | Same flow as the wireframe game, now drawing sprites |
| `ship.entity.go`, `rock.entity.go`, `bullet.entity.go` | Entities; ship and rock draw a texture |
| `save.go` | The high score store |
| `assets/` | `ship.png`, `asteroid.png`, `space-background.png` and the prompts that generated them |

## How assets work

### 1. Embed the files

```go
//go:embed assets/ship.png
var shipPNG []byte
```

`go:embed` copies the file into the program at build time and puts its bytes in a variable. The game does not need the `assets/` folder at run time.

### 2. List them in a manifest

```go
engine.WithAssets(assets.Manifest{
    Textures: []assets.TextureSpec{
        {Name: texShip, Path: "ship.png", Data: shipPNG, Smooth: true},
        ...
    },
})
```

- **Name** is what you call it later. `texShip` is a constant for the string `"ship"`, so a typo is a compile error.
- **Data** is the embedded bytes. `Path` is only used for its extension, which tells raylib the format.
- **Smooth** turns on mipmaps. The images are 1254 px wide but drawn at about 100 px, and without it they look rough.

### 3. The engine loads them

You call no load function. `Run` opens the window, loads the manifest, then enters the first scene. A file that cannot be decoded makes `Run` return an error naming it.

### 4. Fetch by name, keep the texture

`newShip` calls `ctx.Assets.Texture(texShip)` once and stores it. The ship never looks it up per frame. Asking for a name that was never loaded panics with that name.

### 5. Draw it (`draw.go`)

- `drawBackground(ctx)` stretches the background over the whole window.
- `drawSprite(tex, pos, size, rotation)` draws a texture centered on `pos` at `size × size` pixels, rotated around its center. It works for the ship and the rocks because both images are square with the object in the middle.

Rocks also get a random starting angle and a spin speed, so they tumble as they drift.

## Collisions are circles, not pictures

The game never checks whether the images overlap. It uses two circles per pair, sized to the visible shape, not the sprite:

| Constant | Value | Meaning |
|---|---|---|
| `shipSize` / `shipCollisionRadius` | 90 / 22 | sprite size including transparent padding / hit circle |
| `rockSize` / `rockRadius` | 110 / 38 | same for rocks |

If hits feel too easy or too hard, change the radius, not the sprite size.

## How display settings work

The menu shows `F fullscreen`, `V vsync` and `R resolution`.

1. `settings.go` defines `prefs`, which embeds `engine.Display`, and stores it in `~/.config/asteroids_2d/settings.json`.
2. `main.go` reads it and starts the window with `engine.WithDisplay(loadPrefs().Display)`.
3. On a key press, `menu.scene.go` calls `changeDisplay`, which copies `ctx.Display()`, edits it, calls `ctx.ApplyDisplay(d)` so the open window changes immediately, and saves the result.

The next run starts with the saved window. Scenes lay themselves out from `rl.GetScreenWidth()` and `GetScreenHeight()`, so they adapt to any size. Sprites stay the same pixel size, so a very large window shows a small ship.

## Things to try

- Change `rockSize` from `110` to `160`, then raise `rockRadius` to match.
- Add a second texture, for example a bullet sprite: embed it, add it to the manifest, fetch it in `newBullet`, and draw it with `drawSprite`.
- Add a settings entry for `FPS`, using the same `changeDisplay` pattern.
