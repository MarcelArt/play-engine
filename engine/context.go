package engine

import (
	"github.com/MarcelArt/play-engine/assets"
	"github.com/MarcelArt/play-engine/input"
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

	quit bool
}

// Quit ends the game after the current frame.
func (c *Context) Quit() { c.quit = true }
