package engine

// Scene is one screen of a game (menu, level, pause overlay...).
//
// Scenes never import each other. They ask Context.Scenes to move to another
// scene by name, and the factories are registered once in main.
type Scene interface {
	// Enter runs once, after the window exists, when the scene becomes part of
	// the stack. Load assets and build state here, not in a constructor.
	Enter(ctx *Context)
	// Update runs every frame, only for the top scene.
	Update(ctx *Context, dt float32)
	// FixedUpdate runs zero or more times per frame at Config.FixedStep, only
	// for the top scene. Put simulation and physics here.
	FixedUpdate(ctx *Context, dt float32)
	// Draw runs every frame for every scene in the stack, bottom first, so
	// overlays such as a pause menu draw on top of the scene beneath.
	Draw(ctx *Context)
	// Exit runs once when the scene leaves the stack.
	Exit(ctx *Context)
}

// BaseScene provides no-op implementations; embed it and override what you use.
type BaseScene struct{}

func (BaseScene) Enter(*Context)                {}
func (BaseScene) Update(*Context, float32)      {}
func (BaseScene) FixedUpdate(*Context, float32) {}
func (BaseScene) Draw(*Context)                 {}
func (BaseScene) Exit(*Context)                 {}

// Factory builds a scene. It runs when the scene is switched or pushed to, so
// every visit gets fresh state.
type Factory func(ctx *Context) Scene
