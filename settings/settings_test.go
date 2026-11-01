package settings

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

type prefs struct {
	Volume     float32
	Fullscreen bool
	Name       string
}

var defaults = prefs{Volume: 0.8, Name: "player"}

func newStore(t *testing.T) *Store[prefs] {
	return NewAt(filepath.Join(t.TempDir(), "nested", "settings.json"), defaults)
}

func TestMissingFileGivesDefaults(t *testing.T) {
	s := newStore(t)
	got, err := s.Load()
	if !errors.Is(err, ErrNotFound) || got != defaults {
		t.Fatalf("Load = %+v, %v", got, err)
	}
}

func TestRoundTrip(t *testing.T) {
	s := newStore(t)
	want := prefs{Volume: 0.2, Fullscreen: true, Name: "ada"}
	if err := s.Save(want); err != nil {
		t.Fatal(err)
	}
	got, err := s.Load()
	if err != nil || got != want {
		t.Fatalf("Load = %+v, %v", got, err)
	}
}

func TestPartialFileKeepsDefaultsForAbsentKeys(t *testing.T) {
	s := newStore(t)
	os.MkdirAll(filepath.Dir(s.Path()), 0o755)
	os.WriteFile(s.Path(), []byte(`{"Fullscreen": true}`), 0o644)

	got, err := s.Load()
	if err != nil {
		t.Fatal(err)
	}
	if !got.Fullscreen || got.Volume != 0.8 || got.Name != "player" {
		t.Fatalf("got %+v", got)
	}
}

func TestCorruptFileGivesDefaultsAndError(t *testing.T) {
	s := newStore(t)
	os.MkdirAll(filepath.Dir(s.Path()), 0o755)
	os.WriteFile(s.Path(), []byte("{nope"), 0o644)

	got, err := s.Load()
	if err == nil || errors.Is(err, ErrNotFound) || got != defaults {
		t.Fatalf("Load = %+v, %v", got, err)
	}
	if s.LoadOrDefaults() != defaults {
		t.Fatal("LoadOrDefaults should fall back")
	}
}

func TestReset(t *testing.T) {
	s := newStore(t)
	s.Save(prefs{Volume: 0.1})
	got, err := s.Reset()
	if err != nil || got != defaults {
		t.Fatalf("Reset = %+v, %v", got, err)
	}
	if loaded, _ := s.Load(); loaded != defaults {
		t.Fatalf("after Reset Load = %+v", loaded)
	}
}
