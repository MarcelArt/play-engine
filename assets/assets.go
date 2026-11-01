// Package assets loads and owns textures and fonts. A Manager is created per
// game, so there is no shared global registry.
package assets

import (
	"fmt"
	"path/filepath"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// TextureSpec describes one texture to load.
type TextureSpec struct {
	Name string
	// Path is read from disk relative to the working directory. When Data is
	// set it is not read, but its extension still tells raylib the file type.
	Path string
	// Data holds the file's bytes, typically from go:embed, so the game does
	// not depend on the working directory.
	Data []byte
	// Smooth selects mipmapped trilinear filtering, which keeps large art
	// clean when drawn small. The default is point filtering, so pixel art
	// stays crisp.
	Smooth bool
}

// FontSpec describes one font to load.
type FontSpec struct {
	Name string
	Path string
}

// Manifest lists everything to load at startup.
type Manifest struct {
	Textures []TextureSpec
	Fonts    []FontSpec
}

// Manager holds loaded assets by name.
type Manager struct {
	textures map[string]rl.Texture2D
	fonts    map[string]rl.Font
}

func New() *Manager {
	return &Manager{
		textures: make(map[string]rl.Texture2D),
		fonts:    make(map[string]rl.Font),
	}
}

// Load loads every entry in m. It needs an open window. On error, whatever
// loaded so far stays owned by the Manager and is freed by UnloadAll.
func (a *Manager) Load(m Manifest) error {
	for _, s := range m.Textures {
		if err := a.LoadTexture(s); err != nil {
			return err
		}
	}
	for _, s := range m.Fonts {
		if err := a.LoadFont(s); err != nil {
			return err
		}
	}
	return nil
}

// LoadTexture loads one texture, replacing any with the same name.
func (a *Manager) LoadTexture(s TextureSpec) error {
	var tex rl.Texture2D
	if len(s.Data) > 0 {
		img := rl.LoadImageFromMemory(filepath.Ext(s.Path), s.Data, int32(len(s.Data)))
		tex = rl.LoadTextureFromImage(img)
		rl.UnloadImage(img)
	} else {
		tex = rl.LoadTexture(s.Path)
	}
	if tex.ID == 0 {
		return fmt.Errorf("texture %q: cannot load %s", s.Name, s.Path)
	}
	if s.Smooth {
		rl.GenTextureMipmaps(&tex)
		rl.SetTextureFilter(tex, rl.FilterTrilinear)
	} else {
		rl.SetTextureFilter(tex, rl.FilterPoint)
	}
	if old, ok := a.textures[s.Name]; ok {
		rl.UnloadTexture(old)
	}
	a.textures[s.Name] = tex
	return nil
}

// LoadFont loads one font, replacing any with the same name.
func (a *Manager) LoadFont(s FontSpec) error {
	f := rl.LoadFont(s.Path)
	if f.Texture.ID == 0 {
		return fmt.Errorf("font %q: cannot load %s", s.Name, s.Path)
	}
	if old, ok := a.fonts[s.Name]; ok {
		rl.UnloadFont(old)
	}
	a.fonts[s.Name] = f
	return nil
}

// Texture returns the named texture. A missing name is a programming error
// (a typo would otherwise draw nothing, silently), so it panics with the name.
func (a *Manager) Texture(name string) rl.Texture2D {
	t, ok := a.textures[name]
	if !ok {
		panic(fmt.Sprintf("assets: texture %q not loaded", name))
	}
	return t
}

// Font returns the named font; it panics like Texture when missing.
func (a *Manager) Font(name string) rl.Font {
	f, ok := a.fonts[name]
	if !ok {
		panic(fmt.Sprintf("assets: font %q not loaded", name))
	}
	return f
}

// UnloadAll frees everything. It is safe to call twice.
func (a *Manager) UnloadAll() {
	for k, t := range a.textures {
		rl.UnloadTexture(t)
		delete(a.textures, k)
	}
	for k, f := range a.fonts {
		rl.UnloadFont(f)
		delete(a.fonts, k)
	}
}
