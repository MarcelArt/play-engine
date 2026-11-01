package main

import (
	"log"

	"github.com/MarcelArt/play-engine/engine"
	"github.com/MarcelArt/play-engine/settings"
)

const appName = "asteroids_2d"

// prefs is what gets saved to ~/.config/asteroids_2d/settings.json. Embedding
// engine.Display puts width, height, fullscreen, vsync and fps at the top
// level of the file.
type prefs struct {
	engine.Display
}

var defaultPrefs = prefs{Width: 800, Height: 540, FPS: 60}

// prefsStore returns nil when the config directory is unavailable; the game
// still runs, it just cannot remember settings.
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

// changeDisplay applies the edit to the open window and saves the result.
func changeDisplay(ctx *engine.Context, edit func(*engine.Display)) {
	d := ctx.Display()
	edit(&d)
	ctx.ApplyDisplay(d)

	if s := prefsStore(); s != nil {
		if err := s.Save(prefs{Display: ctx.Display()}); err != nil {
			log.Println("could not save settings:", err)
		}
	}
}
