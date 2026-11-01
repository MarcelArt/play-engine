package engine

import (
	"github.com/MarcelArt/play-engine/assets"
	"github.com/MarcelArt/play-engine/input"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Context is the set of services handed to every scene. It replaces package
// globals: a scene reaches assets, input and scene changes only through it,
// and two games never share state.
type Context struct {
	// Name is the game's title; use it to namespace saves (save.New(ctx.Name, ...)).
	Name   string
	Assets *assets.Manager
	Input  *input.Map
	Scenes *SceneManager

	quit    bool
	display Display
}

// Quit ends the game after the current frame.
func (c *Context) Quit() { c.quit = true }

// SetMouseCaptured hides the cursor and locks it to the window, so mouse
// movement keeps producing deltas (first-person look). Release it in menus
// and pause screens, typically in a scene's Enter and Exit.
func (c *Context) SetMouseCaptured(captured bool) {
	if captured {
		rl.DisableCursor()
	} else {
		rl.EnableCursor()
	}
}

// MouseCaptured reports whether the cursor is currently captured.
func (c *Context) MouseCaptured() bool { return rl.IsCursorHidden() }
