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
