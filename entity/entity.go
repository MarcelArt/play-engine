// Package entity is an optional helper for scenes that keep a list of
// things to update and draw. Games that drive their own data (a tile grid)
// can ignore it.
package entity

// Entity is anything a World can hold.
type Entity interface {
	Update(dt float32)
	Draw()
	// Alive reports whether the entity should stay in the world.
	Alive() bool
}

// Lifecycle can be embedded to give an entity Alive and Kill.
type Lifecycle struct{ dead bool }

func (l *Lifecycle) Alive() bool { return !l.dead }
func (l *Lifecycle) Kill()       { l.dead = true }

// World is an ordered list of entities. Entities added during Update are first
// updated on the next call, and dead ones are swept at the end of Update, so
// callers never remove while iterating.
type World[T Entity] struct {
	items []T
}

func (w *World[T]) Add(e T)  { w.items = append(w.items, e) }
func (w *World[T]) Len() int { return len(w.items) }
func (w *World[T]) Clear()   { clear(w.items); w.items = w.items[:0] }

// Items returns the live slice; do not retain it across Update.
func (w *World[T]) Items() []T { return w.items }

func (w *World[T]) Update(dt float32) {
	for _, e := range w.items[:len(w.items):len(w.items)] {
		if e.Alive() {
			e.Update(dt)
		}
	}
	w.Sweep()
}

func (w *World[T]) Draw() {
	for _, e := range w.items {
		if e.Alive() {
			e.Draw()
		}
	}
}

// Sweep removes dead entities. Call it after collision handling if you kill
// entities outside Update.
func (w *World[T]) Sweep() {
	kept := w.items[:0]
	for _, e := range w.items {
		if e.Alive() {
			kept = append(kept, e)
		}
	}
	clear(w.items[len(kept):])
	w.items = kept
}
