package main

import (
	"azul3d.org/engine/gfx/window"
	"common"
	"games"
	"log"
)

func main() {
	// set the property of window
	probs := window.NewProps()
	probs.SetSize(common.WINDOW_WIDTH, common.WINDOW_HEIGHT)
	// init games
	gm_env := &games.Game{}
	err := gm_env.Init()
	if nil != err {
		log.Fatal(err)
	}
	// window run
	window.Run(gm_env.GfxLoop, probs)
}
