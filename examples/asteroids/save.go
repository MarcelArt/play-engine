package main

import (
	"log"

	"github.com/MarcelArt/play-engine/engine"
	"github.com/MarcelArt/play-engine/save"
)

type saveFile struct {
	HighScore uint
}

func highScoreStore(ctx *engine.Context) *save.Store[saveFile] {
	s, err := save.New[saveFile](ctx.Name, "save.json")
	if err != nil {
		log.Println("save disabled:", err)
		return nil
	}
	return s
}
