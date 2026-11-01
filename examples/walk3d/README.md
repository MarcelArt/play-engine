# walk3d

A first-person 3D walkaround: mouse look, WASD movement, sprinting, jumping, simple collision with blocks, a pause scene that frees the mouse, and saved display settings. Read it to see how to build a 3D game with play-engine.

```bash
go run ./examples/walk3d      # from the repository root
```

| Key | Action |
|---|---|
| Mouse (or right stick) | look |
| W A S D or arrows (or left stick) | move |
| Shift | sprint |
| Space (or gamepad bottom face button) | jump |
| F | toggle fullscreen |
| P | pause (frees the mouse) / resume |
| Esc | quit |

## What it shows

- `camera.FirstPerson` for mouse look and flat movement vectors
- Capturing and releasing the mouse with `ctx.SetMouseCaptured`
- Gamepad sticks with `ctx.Input.Stick`
- A pause scene that changes the mouse state on `Enter` and `Exit`
- `engine.WithMSAA()` for smooth edges
- Display, sensitivity and field of view saved with `settings`

## Files

| File | Role |
|---|---|
| `main.go` | Builds the game, binds actions, registers the scenes |
| `settings.go` | The saved settings: display, sensitivity, FOV |
| `play.scene.go` | The scene: owns the world and player, draws 3D then the HUD |
| `pause.scene.go` | Overlay that frees the mouse while open |
| `player.go` | The first-person controller: look, move, jump, gravity |
| `world.go` | The floor and blocks, plus push-out and ground checks |

## Coordinates

raylib is **Y-up and right-handed**. Yaw 0 looks down -Z, positive yaw turns left, and one unit is one metre. The player's position is stored as the position of their **feet**, and the camera sits `eyeHeight` (1.7) above it.

## How it runs

### 1. `main.go`

- Loads the saved settings first, so the window opens at the saved size.
- `WithDisplay` sets size, fullscreen, vsync and FPS from them.
- `WithMSAA` smooths polygon edges, and `WithBackground` makes the sky.
- Registers `play` and `pause`, and starts on `play`.

### 2. The play scene starts (`play.scene.go`)

`Enter` builds the world and the player, then calls `ctx.SetMouseCaptured(true)`. That hides the cursor and locks it to the window, so mouse movement keeps producing deltas. `Exit` releases it, so the cursor is never left locked.

### 3. Every frame (`player.go`)

`Update` runs the controller. It moves in `Update`, not `FixedUpdate`, because it follows input directly and no interpolation is needed.

**Look.** While the mouse is captured, `ctx.Input.MouseDelta()` is passed to `cam.Look`, which turns the camera and clamps the pitch. The first frame after capturing is skipped, because the cursor jumps to the window center and reports a large delta. The right stick turns at a fixed rate multiplied by `dt`.

**Move.**
1. Keys and the left stick both produce a forward value and a strafe value in -1..1.
2. These scale `cam.Forward()` and `cam.Right()`, which are flat on the ground. Walking forward never climbs or dives with the pitch.
3. A diagonal is normalised so it is not faster than straight movement.
4. Shift multiplies the speed.

**Collide sideways.** After moving, `world.pushOut` moves the player out of any block they overlap. Blocks whose top is within `stepHeight` (0.55 m) of the feet are ignored, so small steps are walked up instead of blocking.

**Jump and land.** Gravity pulls the player down each frame. Space sets an upward speed when grounded. Then `world.groundAt` finds the highest surface under the player, either a block top or the floor, and the player lands when they reach it. The jump reaches about 1.1 m, so most blocks can be jumped onto and the tallest cannot.

Finally `cam.SetEye` moves the camera to the player's eye.

### 4. Drawing

`Draw` runs the 3D pass between `cam.Begin()` and `cam.End()`: the floor, a grid, and each block as a coloured cube with dark edges. The HUD (crosshair, controls, position) is drawn after `End`, in screen space.

### 5. Pause (`pause.scene.go`)

P pushes the pause scene. Its `Enter` releases the mouse so you can use the cursor. The play scene stops updating but keeps drawing under the dim overlay. P pops the pause scene, and its `Exit` captures the mouse again.

### 6. The world (`world.go`)

- **Stairs** of 0.5 m steps sit to the right of the spawn, for trying out step-up and jumping.
- **Random blocks** come from a fixed seed, so the layout is the same every run. None spawn within 5 m of the start.
- Blocks are axis-aligned boxes. This is the example's own simple collision, written here and not in the library.

### 7. Settings (`settings.go`)

The saved file (`~/.config/walk3d/settings.json`) holds the display fields plus `sensitivity` (radians per pixel) and `fov` (degrees). F toggles fullscreen and saves it. Edit the file to change sensitivity or FOV; there is no in-game menu for them.

## Limits

- Collision is a vertical cylinder against axis-aligned boxes. There are no slopes, no walls taller than the world, and no collision with the world's edge: the floor simply ends.
- Movement runs per frame, so jump height varies slightly with the frame rate.

## Things to try

- Change `jumpSpeed` or `gravity` in `player.go` to change how high and floaty jumps feel.
- Lower `sensitivity` in the settings file if the look feels too fast.
- Add more stairs in `newWorld`, or make a block taller than the jump so it forces a detour.
- Add a wall around the floor by adding four long thin boxes at `±worldHalf`.
