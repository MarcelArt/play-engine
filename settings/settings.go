// Package settings persists user preferences (volume, key bindings, window
// options) as JSON in the user's config directory, apart from save data.
package settings

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/MarcelArt/play-engine/save"
	"github.com/adrg/xdg"
)

// ErrNotFound is returned by Load when no settings file exists yet. The value
// returned with it is the defaults, so it is safe to use as is.
var ErrNotFound = errors.New("settings: no settings file")

// Store reads and writes one settings value of type T. It always has defaults:
// keys missing from the file keep their default value, so adding a new
// setting in a later version never resets the old ones.
type Store[T any] struct {
	path     string
	defaults T
}

// New returns a Store at <user config dir>/<app>/<file> (~/.config/<app>/<file>
// on Linux). Namespacing by app keeps games from sharing settings.
func New[T any](app, file string, defaults T) (*Store[T], error) {
	p, err := xdg.ConfigFile(filepath.Join(app, file))
	if err != nil {
		return nil, err
	}
	return &Store[T]{path: p, defaults: defaults}, nil
}

// NewAt returns a Store at an explicit path (useful for tests).
func NewAt[T any](path string, defaults T) *Store[T] {
	return &Store[T]{path: path, defaults: defaults}
}

// Path is where the settings are stored.
func (s *Store[T]) Path() string { return s.path }

// Defaults returns the default settings.
func (s *Store[T]) Defaults() T { return s.defaults }

// Load returns the stored settings laid over the defaults. The returned value
// is always usable: when the file is missing it is the defaults and the error
// is ErrNotFound; when the file is corrupt it is the defaults and the error
// says why. Slices and maps in the file replace or merge per encoding/json
// rules, so prefer plain fields for settings.
func (s *Store[T]) Load() (T, error) {
	v := s.defaults
	b, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return s.defaults, ErrNotFound
	}
	if err != nil {
		return s.defaults, err
	}
	if err := json.Unmarshal(b, &v); err != nil {
		return s.defaults, err
	}
	return v, nil
}

// LoadOrDefaults is Load for callers that do not care why defaults were used.
func (s *Store[T]) LoadOrDefaults() T {
	v, _ := s.Load()
	return v
}

// Save writes v atomically, so a crash never leaves a half-written file.
func (s *Store[T]) Save(v T) error {
	return save.NewAt[T](s.path).Save(v)
}

// Reset writes the defaults over the stored settings and returns them.
func (s *Store[T]) Reset() (T, error) {
	return s.defaults, s.Save(s.defaults)
}
