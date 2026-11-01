package engine

import rl "github.com/gen2brain/raylib-go/raylib"

// Display is the window setup a player can change: size, fullscreen, vsync
// and frame cap. It has JSON tags so it can be embedded in a settings struct
// and persisted with the settings package:
//
//	type Prefs struct {
//		engine.Display
//		Volume float32
//	}
type Display struct {
	Width      int32 `json:"width"`
	Height     int32 `json:"height"`
	Fullscreen bool  `json:"fullscreen"`
	VSync      bool  `json:"vsync"`
	// FPS caps the frame rate; 0 means uncapped.
	FPS int32 `json:"fps"`
}

const (
	minDisplayWidth  = 320
	minDisplayHeight = 200
)

// DefaultDisplay is 1280×720 windowed at 60 FPS.
func DefaultDisplay() Display {
	return Display{Width: 1280, Height: 720, FPS: 60}
}

// Normalize returns d with values a window can actually use: size at least
// 320×200 and a non-negative frame cap. Settings read from a hand-edited file
// pass through it, so a typo cannot create a 0×0 window.
func (d Display) Normalize() Display {
	d.Width = max(d.Width, minDisplayWidth)
	d.Height = max(d.Height, minDisplayHeight)
	d.FPS = max(d.FPS, 0)
	return d
}

// Resolution is a window size.
type Resolution struct{ Width, Height int32 }

// CommonResolutions are sizes a settings menu can offer, smallest first.
var CommonResolutions = []Resolution{
	{800, 600},
	{1024, 768},
	{1280, 720},
	{1600, 900},
	{1920, 1080},
	{2560, 1440},
}

// NextResolution returns the entry after w×h in CommonResolutions, wrapping
// around. A size not in the list returns the first entry.
func NextResolution(w, h int32) Resolution {
	for i, r := range CommonResolutions {
		if r.Width == w && r.Height == h {
			return CommonResolutions[(i+1)%len(CommonResolutions)]
		}
	}
	return CommonResolutions[0]
}

// Display returns the display settings last applied (or the initial ones).
func (c *Context) Display() Display { return c.display }

// ApplyDisplay changes the open window to match d. It takes effect
// immediately, so call it from a scene's Update (for example when a settings
// menu changes a value), then save d with the settings package.
//
// Leaving fullscreen restores the windowed size and centers the window.
// Scenes that read rl.GetScreenWidth/GetScreenHeight adapt on their own.
func (c *Context) ApplyDisplay(d Display) {
	d = d.Normalize()

	if d.VSync {
		rl.SetWindowState(rl.FlagVsyncHint)
	} else {
		rl.ClearWindowState(rl.FlagVsyncHint)
	}
	rl.SetTargetFPS(d.FPS)

	switch {
	case d.Fullscreen && !rl.IsWindowFullscreen():
		rl.ToggleFullscreen()
	case !d.Fullscreen && rl.IsWindowFullscreen():
		rl.ToggleFullscreen()
	}

	if !d.Fullscreen {
		rl.SetWindowSize(int(d.Width), int(d.Height))
		centerWindow(int(d.Width), int(d.Height))
	}

	c.display = d
}

func centerWindow(w, h int) {
	m := rl.GetCurrentMonitor()
	mw, mh := rl.GetMonitorWidth(m), rl.GetMonitorHeight(m)
	rl.SetWindowPosition(max((mw-w)/2, 0), max((mh-h)/2, 0))
}
