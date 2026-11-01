# asteroids

The smallest complete play-engine game: a wireframe Asteroids with a menu, a play scene, a pause overlay, and a saved high score. It uses no asset files, so it is the best place to start reading.

```bash
go run ./examples/asteroids      # from the repository root
```

| Key | Action |
|---|---|
| Enter | start / retry |
| Arrows or A/D | rotate |
| Up or W | thrust |
| Space | shoot |
| P | pause / resume |
| Esc | quit |

## What it shows

- Registering scenes by name and moving between them (`Switch`, `Push`, `Pop`)
- Input by action name instead of key
- `Update` for input and `FixedUpdate` for the simulation
- `entity.World` for lists of bullets and rocks
- A saved high score with the `save` package

## Files

| File | Role |
|---|---|
| `main.go` | Builds the game, binds input actions, registers the three scenes, runs the menu |
| `menu.scene.go` | Title screen; Enter switches to play |
| `play.scene.go` | The game: owns the ship, bullets, rocks, score and game-over state |
| `pause.scene.go` | Overlay pushed on top of play |
| `ship.entity.go` | The player ship: rotation, thrust, drag |
| `bullet.entity.go` | A bullet that flies straight and expires |
| `rock.entity.go` | A drifting asteroid with a jagged outline |
| `save.go` | The save file type and its store |

## How it runs

### 1. `main.go` sets up the game

`engine.New("asteroids", ...)` takes the window size and the key bindings. The bindings map **actions** such as `actionShoot` to keys such as Space, so the rest of the code never names a key. Then `g.Register(name, factory)` adds the three scenes, and `g.Run(sceneMenu)` opens the window on the menu.

### 2. The menu

When the menu is created it loads the high score. Each frame it draws the title; pressing Enter calls `ctx.Scenes.Switch(scenePlay)`, which replaces the menu with a new play scene.

### 3. The play scene starts

`newPlayScene` loads the high score and sets the rock timer. `Enter` then puts the ship in the middle of the screen. The scene holds the ship plus two `entity.World` lists, one for bullets and one for rocks.

### 4. Every frame: two kinds of update

**`Update` runs once per frame and handles input:**
- P pushes the pause scene.
- `ship.Update` reads the keys, turns and thrusts the ship, and moves it.
- If Space was pressed, the ship does not create a bullet itself. It sets a flag, and the scene sees it through `TakeShot()` and adds a bullet at the ship's nose. The entity asks, the scene acts.

**`FixedUpdate` runs 60 times a second and handles the world:**
- Move every bullet and rock.
- Count down the spawn timer, then spawn a rock from a random corner aimed at the ship.
- Run `collide()`.

### 5. Collisions

`collide()` uses circles (`rl.CheckCollisionCircles`):

1. Every bullet and rock that touch are both killed, and the score goes up.
2. `Sweep()` removes the killed ones from the lists.
3. If any rock touches the ship, `lose()` runs.

### 6. Pause

`Push(scenePause)` puts the pause scene on top. The play scene stops updating, so everything freezes, but it still draws, so you see the frozen game under the dark overlay. P pops the pause scene, and Enter switches back to the menu.

### 7. Game over and saving

`lose()` sets `dead`. If the score beats the high score it is written with the `save` package (to `~/.local/share/asteroids/save.json` on Linux). While dead, `Update` only waits for Enter, which switches to a fresh play scene.

## Things to try

- Set `rockSpawnInterval` in `play.scene.go` to `0.5` for a harder game.
- Raise `shipThrust` or lower `shipDrag` in `ship.entity.go` for a slipperier ship.
- Add a third scene, for example a settings screen, registered in `main.go` and reached with `Push`.
