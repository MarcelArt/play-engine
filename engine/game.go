// Package engine is a small scene-based game loop on top of raylib.
package engine

import (
	"fmt"

	"github.com/MarcelArt/play-engine/assets"
	"github.com/MarcelArt/play-engine/input"
	"github.com/MarcelArt/play-engine/option"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Game wires a window, the services in Context, and a stack of scenes.
type Game struct {
	title string
	cfg   Config
	ctx   *Context
}

// New creates a game. Nothing touches raylib until Run.
func New(title string, opts ...Option) *Game {
	cfg := defaultConfig()
	option.Apply(&cfg, opts...)

	g := &Game{title: title, cfg: cfg}
	g.ctx = &Context{
		Name:   title,
		Assets: assets.New(),
		Input:  input.New(cfg.Bindings),
	}
	g.ctx.Scenes = newSceneManager(g.ctx)
	return g
}

// Register adds a scene factory under name. Call it before Run.
func (g *Game) Register(name string, f Factory) *Game {
	g.ctx.Scenes.register(name, f)
	return g
}

// Context exposes the services, e.g. for loading extra assets before Run.
func (g *Game) Context() *Context { return g.ctx }

// Run opens the window, enters the start scene and loops until the game quits.
// It returns an error if assets fail to load. Cleanup runs even on panic.
func (g *Game) Run(start string) error {
	if g.cfg.Resizable {
		rl.SetConfigFlags(rl.FlagWindowResizable)
	}
	rl.InitWindow(g.cfg.Width, g.cfg.Height, g.title)
	defer rl.CloseWindow()
	rl.SetExitKey(g.cfg.ExitKey)
	rl.SetTargetFPS(g.cfg.FPS)

	defer g.ctx.Assets.UnloadAll()
	if err := g.ctx.Assets.Load(g.cfg.Manifest); err != nil {
		return fmt.Errorf("engine: loading assets: %w", err)
	}

	scenes := g.ctx.Scenes
	defer scenes.shutdown()
	scenes.Switch(start)
	scenes.apply()

	var acc float32
	for !rl.WindowShouldClose() && !g.ctx.quit {
		dt := rl.GetFrameTime()
		top := scenes.top()

		top.Update(g.ctx, dt)

		// Fixed timestep: carry leftover time into the next frame instead of
		// dropping it, and cap the catch-up so a long stall cannot spiral.
		acc += dt
		for steps := 0; acc >= g.cfg.FixedStep && steps < g.cfg.MaxFixedSteps; steps++ {
			top.FixedUpdate(g.ctx, g.cfg.FixedStep)
			acc -= g.cfg.FixedStep
		}
		if acc >= g.cfg.FixedStep {
			acc = 0
		}

		rl.BeginDrawing()
		rl.ClearBackground(g.cfg.Background)
		for _, s := range scenes.stack {
			s.Draw(g.ctx)
		}
		if g.cfg.ShowFPS {
			rl.DrawFPS(0, 0)
		}
		rl.EndDrawing()

		// Scene changes land between frames, so scenes may request them freely.
		scenes.apply()
		if len(scenes.stack) == 0 {
			break
		}
	}
	return nil
}
