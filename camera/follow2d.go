// Package camera provides ready-made cameras: a smoothed, bounds-clamped
// follow camera for 2D games and a first-person camera for 3D.
package camera

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Follow2D keeps a target centered on screen, optionally smoothed and kept
// inside the world. Wrap drawing in Begin and End.
type Follow2D struct {
	// Camera is the underlying raylib camera; Offset and Target are managed by Update.
	Camera rl.Camera2D
	// Smooth is how fast the camera catches up, in 1/seconds. 0 snaps to the target.
	// Higher is snappier; around 6–10 feels natural.
	Smooth float32
	// Bounds keeps the view inside this world rectangle when set. A world
	// smaller than the screen is centered instead.
	Bounds *rl.Rectangle
}

// NewFollow2D returns a camera at the given zoom (1 is one pixel per unit).
func NewFollow2D(zoom float32) *Follow2D {
	return &Follow2D{Camera: rl.NewCamera2D(rl.Vector2Zero(), rl.Vector2Zero(), 0, zoom)}
}

// SetBounds limits the view to a world rectangle in world units.
func (c *Follow2D) SetBounds(r rl.Rectangle) { c.Bounds = &r }

// Snap jumps straight to target, with no smoothing; use it when a level loads.
func (c *Follow2D) Snap(target rl.Vector2) {
	smooth := c.Smooth
	c.Smooth = 0
	c.Update(target, 0)
	c.Smooth = smooth
}

// Update moves the camera toward target and re-centers on the current window
// size, so resizing the window or changing fullscreen needs no extra code.
func (c *Follow2D) Update(target rl.Vector2, dt float32) {
	sw, sh := float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())
	c.Camera.Offset = rl.NewVector2(sw/2, sh/2)

	c.Camera.Target.X = approach(c.Camera.Target.X, target.X, c.Smooth, dt)
	c.Camera.Target.Y = approach(c.Camera.Target.Y, target.Y, c.Smooth, dt)

	if c.Bounds != nil {
		zoom := max(c.Camera.Zoom, 0.0001)
		c.Camera.Target.X = clampAxis(c.Camera.Target.X, c.Bounds.X, c.Bounds.Width, sw/zoom)
		c.Camera.Target.Y = clampAxis(c.Camera.Target.Y, c.Bounds.Y, c.Bounds.Height, sh/zoom)
	}
}

// Begin starts drawing in world space.
func (c *Follow2D) Begin() { rl.BeginMode2D(c.Camera) }

// End stops drawing in world space; draw HUD after it.
func (c *Follow2D) End() { rl.EndMode2D() }

// ScreenToWorld converts a screen position, such as the mouse, to world space.
func (c *Follow2D) ScreenToWorld(p rl.Vector2) rl.Vector2 {
	return rl.GetScreenToWorld2D(p, c.Camera)
}

// WorldToScreen converts a world position to a screen position.
func (c *Follow2D) WorldToScreen(p rl.Vector2) rl.Vector2 {
	return rl.GetWorldToScreen2D(p, c.Camera)
}

// approach moves cur toward dst. The exponential form makes the result the
// same at any frame rate; rate <= 0 snaps.
func approach(cur, dst, rate, dt float32) float32 {
	if rate <= 0 {
		return dst
	}
	t := 1 - float32(math.Exp(float64(-rate*dt)))
	return cur + (dst-cur)*t
}

// clampAxis keeps a view of viewSize centered on target inside [min, min+size].
// When the world is smaller than the view, it centers on the world.
func clampAxis(target, min, size, viewSize float32) float32 {
	if size <= viewSize {
		return min + size/2
	}
	half := viewSize / 2
	return rl.Clamp(target, min+half, min+size-half)
}
