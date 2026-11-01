package engine

import (
	"github.com/MarcelArt/play-engine/assets"
	"github.com/MarcelArt/play-engine/input"
	"github.com/MarcelArt/play-engine/option"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Config is the window and loop setup. Change it through Options.
type Config struct {
	Width, Height int32
	FPS           int32
	Resizable     bool
	Fullscreen    bool
	VSync         bool
	MSAA          bool
	ShowFPS       bool
	// FixedStep is the FixedUpdate interval in seconds.
	FixedStep float32
	// MaxFixedSteps caps catch-up steps per frame so a stall cannot spiral.
	MaxFixedSteps int
	Background    rl.Color
	// ExitKey closes the game; 0 disables it (raylib's default is Escape).
	ExitKey  int32
	Manifest assets.Manifest
	Bindings map[input.Action]input.Binding
}

// Option configures a Game.
type Option = option.Func[Config]

func defaultConfig() Config {
	return Config{
		Width:         1280,
		Height:        720,
		FPS:           60,
		FixedStep:     1.0 / 60,
		MaxFixedSteps: 5,
		Background:    rl.Black,
		ExitKey:       rl.KeyEscape,
	}
}

func WithSize(w, h int32) Option { return func(c *Config) { c.Width, c.Height = w, h } }
func WithFPS(fps int32) Option   { return func(c *Config) { c.FPS = fps } }
func WithResizable() Option      { return func(c *Config) { c.Resizable = true } }

// WithMSAA turns on 4× multisampling, which smooths polygon edges in 3D.
func WithMSAA() Option { return func(c *Config) { c.MSAA = true } }

func WithFPSCounter() Option             { return func(c *Config) { c.ShowFPS = true } }
func WithBackground(col rl.Color) Option { return func(c *Config) { c.Background = col } }

// WithDisplay sets the initial window from a Display, typically one loaded
// from settings. It replaces WithSize and WithFPS.
func WithDisplay(d Display) Option {
	return func(c *Config) {
		d = d.Normalize()
		c.Width, c.Height, c.FPS = d.Width, d.Height, d.FPS
		c.Fullscreen, c.VSync = d.Fullscreen, d.VSync
	}
}

// WithExitKey sets the key that quits the game; pass 0 to disable it.
func WithExitKey(key int32) Option { return func(c *Config) { c.ExitKey = key } }

// WithFixedStep sets the FixedUpdate interval in seconds (default 1/60).
func WithFixedStep(seconds float32) Option { return func(c *Config) { c.FixedStep = seconds } }

// WithAssets loads the manifest right after the window opens, before the first scene.
func WithAssets(m assets.Manifest) Option { return func(c *Config) { c.Manifest = m } }

// WithBindings sets the input action bindings.
func WithBindings(b map[input.Action]input.Binding) Option {
	return func(c *Config) { c.Bindings = b }
}
