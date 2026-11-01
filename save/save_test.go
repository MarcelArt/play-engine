package save

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

type data struct{ HighScore uint }

func TestRoundTrip(t *testing.T) {
	s := NewAt[data](filepath.Join(t.TempDir(), "nested", "save.json"))

	if _, err := s.Load(); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Load before Save: %v, want ErrNotFound", err)
	}
	if got := s.LoadOr(data{HighScore: 3}); got.HighScore != 3 {
		t.Fatalf("LoadOr default = %+v", got)
	}

	if err := s.Save(data{HighScore: 42}); err != nil {
		t.Fatal(err)
	}
	got, err := s.Load()
	if err != nil || got.HighScore != 42 {
		t.Fatalf("Load = %+v, %v", got, err)
	}
}

func TestCorruptFile(t *testing.T) {
	p := filepath.Join(t.TempDir(), "save.json")
	os.WriteFile(p, []byte("{nope"), 0o644)
	s := NewAt[data](p)
	if _, err := s.Load(); err == nil || errors.Is(err, ErrNotFound) {
		t.Fatalf("corrupt file should be a real error, got %v", err)
	}
	if got := s.LoadOr(data{HighScore: 1}); got.HighScore != 1 {
		t.Fatal("LoadOr should fall back on corrupt file")
	}
}
