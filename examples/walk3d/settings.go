package main

import (
	"log"

	"github.com/MarcelArt/play-engine/engine"
	"github.com/MarcelArt/play-engine/settings"
)

const appName = "walk3d"

// prefs is saved to ~/.config/walk3d/settings.json.
type prefs struct {
	engine.Display
	// Sensitivity is radians turned per pixel of mouse movement.
	Sensitivity float32 `json:"sensitivity"`
	// FOV is the vertical field of view in degrees.
	FOV float32 `json:"fov"`
}

var defaultPrefs = prefs{
	Width: 1280, Height: 720, FPS: 0, VSync: true,
	Sensitivity: 0.0025,
	FOV:         70,
}

func prefsStore() *settings.Store[prefs] {
	s, err := settings.New(appName, "settings.json", defaultPrefs)
	if err != nil {
		log.Println("settings disabled:", err)
		return nil
	}
	return s
}

func loadPrefs() prefs {
	if s := prefsStore(); s != nil {
		return s.LoadOrDefaults()
	}
	return defaultPrefs
}

// savePrefs writes p; the display part always reflects the open window.
func savePrefs(ctx *engine.Context, p prefs) {
	p.Display = ctx.Display()
	if s := prefsStore(); s != nil {
		if err := s.Save(p); err != nil {
			log.Println("could not save settings:", err)
		}
	}
}
