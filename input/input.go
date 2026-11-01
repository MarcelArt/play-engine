// Package input maps named actions to physical keys, buttons and gamepad
// buttons, so game code asks "is Shoot pressed" and never names a key.
package input

import (
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Action is a game-level intent such as "jump" or "pause".
type Action string

// Binding lists every physical input that triggers an action.
type Binding struct {
	Keys           []int32
	MouseButtons   []rl.MouseButton
	GamepadButtons []int32 // read from gamepad 0
}

// Map resolves actions to their bindings. Rebinding at runtime is Bind.
type Map struct {
	bindings map[Action]Binding
}

func New(b map[Action]Binding) *Map {
	m := &Map{bindings: make(map[Action]Binding, len(b))}
	for a, bind := range b {
		m.bindings[a] = bind
	}
	return m
}

// Bind replaces the binding for a.
func (m *Map) Bind(a Action, b Binding) { m.bindings[a] = b }

// Keys is shorthand for a keyboard-only Binding.
func Keys(keys ...int32) Binding { return Binding{Keys: keys} }

const gamepad = 0

// Down reports whether any input bound to a is held.
func (m *Map) Down(a Action) bool {
	b := m.bindings[a]
	if slices.ContainsFunc(b.Keys, rl.IsKeyDown) {
		return true
	}
	if slices.ContainsFunc(b.MouseButtons, rl.IsMouseButtonDown) {
		return true
	}
	for _, g := range b.GamepadButtons {
		if rl.IsGamepadButtonDown(gamepad, g) {
			return true
		}
	}
	return false
}

// Pressed reports whether any input bound to a went down this frame.
func (m *Map) Pressed(a Action) bool {
	b := m.bindings[a]
	if slices.ContainsFunc(b.Keys, rl.IsKeyPressed) {
		return true
	}
	if slices.ContainsFunc(b.MouseButtons, rl.IsMouseButtonPressed) {
		return true
	}
	for _, g := range b.GamepadButtons {
		if rl.IsGamepadButtonPressed(gamepad, g) {
			return true
		}
	}
	return false
}

// Released reports whether any input bound to a went up this frame.
func (m *Map) Released(a Action) bool {
	b := m.bindings[a]
	if slices.ContainsFunc(b.Keys, rl.IsKeyReleased) {
		return true
	}
	if slices.ContainsFunc(b.MouseButtons, rl.IsMouseButtonReleased) {
		return true
	}
	for _, g := range b.GamepadButtons {
		if rl.IsGamepadButtonReleased(gamepad, g) {
			return true
		}
	}
	return false
}

// Axis returns -1, 0 or 1 from a negative and positive action; both held cancel out.
func (m *Map) Axis(neg, pos Action) float32 {
	var v float32
	if m.Down(neg) {
		v--
	}
	if m.Down(pos) {
		v++
	}
	return v
}

// Stick selects a gamepad thumbstick.
type Stick uint8

const (
	LeftStick Stick = iota
	RightStick
)

// stickDeadzone ignores the small drift sticks report when released.
const stickDeadzone = 0.15

// Stick returns the thumbstick position on gamepad 0, each axis in [-1, 1],
// with the dead zone removed. It is zero when no gamepad is connected.
func (m *Map) Stick(s Stick) rl.Vector2 {
	if !rl.IsGamepadAvailable(gamepad) {
		return rl.Vector2Zero()
	}
	ax, ay := int32(rl.GamepadAxisLeftX), int32(rl.GamepadAxisLeftY)
	if s == RightStick {
		ax, ay = int32(rl.GamepadAxisRightX), int32(rl.GamepadAxisRightY)
	}
	return rl.NewVector2(
		deadzone(rl.GetGamepadAxisMovement(gamepad, ax)),
		deadzone(rl.GetGamepadAxisMovement(gamepad, ay)),
	)
}

func deadzone(v float32) float32 {
	if v > -stickDeadzone && v < stickDeadzone {
		return 0
	}
	return v
}

// MouseDelta returns how far the mouse moved this frame in pixels. Pair it
// with Context.SetMouseCaptured for first-person look.
func (m *Map) MouseDelta() rl.Vector2 { return rl.GetMouseDelta() }
