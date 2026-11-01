// Package save persists one typed value as JSON in the user's data directory.
package save

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

// ErrNotFound is returned by Load when nothing has been saved yet.
var ErrNotFound = errors.New("save: no save file")

// Store reads and writes a single T at a fixed path.
type Store[T any] struct {
	path string
}

// New returns a Store at <user data dir>/<app>/<file>. Namespacing by app keeps
// different games from overwriting each other's saves.
func New[T any](app, file string) (*Store[T], error) {
	p, err := xdg.DataFile(filepath.Join(app, file))
	if err != nil {
		return nil, err
	}
	return &Store[T]{path: p}, nil
}

// NewAt returns a Store at an explicit path (useful for tests).
func NewAt[T any](path string) *Store[T] {
	return &Store[T]{path: path}
}

// Path is where the value is stored.
func (s *Store[T]) Path() string { return s.path }

// Load returns the saved value, or ErrNotFound when none exists.
func (s *Store[T]) Load() (T, error) {
	var v T
	b, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return v, ErrNotFound
	}
	if err != nil {
		return v, err
	}
	err = json.Unmarshal(b, &v)
	return v, err
}

// LoadOr returns the saved value, or def when none exists or it is unreadable.
func (s *Store[T]) LoadOr(def T) T {
	v, err := s.Load()
	if err != nil {
		return def
	}
	return v
}

// Save writes v atomically, so a crash never leaves a half-written file.
func (s *Store[T]) Save(v T) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
