package engine

import "fmt"

type sceneOp struct {
	kind sceneOpKind
	name string
}

type sceneOpKind uint8

const (
	opSwitch sceneOpKind = iota
	opPush
	opPop
)

// SceneManager owns the scene stack. Changes requested during a frame are
// queued and applied between frames, so a scene can safely ask for a change
// from inside its own Update.
type SceneManager struct {
	ctx       *Context
	factories map[string]Factory
	stack     []Scene
	pending   []sceneOp
}

func newSceneManager(ctx *Context) *SceneManager {
	return &SceneManager{ctx: ctx, factories: make(map[string]Factory)}
}

func (m *SceneManager) register(name string, f Factory) {
	if _, dup := m.factories[name]; dup {
		panic(fmt.Sprintf("engine: scene %q registered twice", name))
	}
	m.factories[name] = f
}

// Switch replaces the whole stack with the named scene.
func (m *SceneManager) Switch(name string) { m.pending = append(m.pending, sceneOp{opSwitch, name}) }

// Push puts the named scene on top; the scene below keeps drawing but stops updating.
func (m *SceneManager) Push(name string) { m.pending = append(m.pending, sceneOp{opPush, name}) }

// Pop removes the top scene. Popping the last scene quits the game.
func (m *SceneManager) Pop() { m.pending = append(m.pending, sceneOp{kind: opPop}) }

func (m *SceneManager) top() Scene {
	if len(m.stack) == 0 {
		return nil
	}
	return m.stack[len(m.stack)-1]
}

func (m *SceneManager) apply() {
	ops := m.pending
	m.pending = nil
	for _, op := range ops {
		switch op.kind {
		case opSwitch:
			for len(m.stack) > 0 {
				m.popTop()
			}
			m.pushNamed(op.name)
		case opPush:
			m.pushNamed(op.name)
		case opPop:
			m.popTop()
		}
	}
	if len(m.stack) == 0 && len(m.pending) == 0 {
		m.ctx.Quit()
	}
}

func (m *SceneManager) pushNamed(name string) {
	f, ok := m.factories[name]
	if !ok {
		panic(fmt.Sprintf("engine: unknown scene %q (register it before Run)", name))
	}
	s := f(m.ctx)
	m.stack = append(m.stack, s)
	s.Enter(m.ctx)
}

func (m *SceneManager) popTop() {
	if len(m.stack) == 0 {
		return
	}
	s := m.stack[len(m.stack)-1]
	m.stack = m.stack[:len(m.stack)-1]
	s.Exit(m.ctx)
}

func (m *SceneManager) shutdown() {
	m.pending = nil
	for len(m.stack) > 0 {
		m.popTop()
	}
}
