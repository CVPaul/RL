package main

import (
	"azul3d.org/engine/gfx/window"
	"common"
	"games"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Print("Error||Lenght of args must larger than 2")
		os.Exit(-1)
	}
	// set the property of window
	probs := window.NewProps()
	probs.SetSize(common.WINDOW_WIDTH, common.WINDOW_HEIGHT)
	// init games
	gm_env := &games.Game{}
	err := gm_env.InitTesting(os.Args[1])
	if nil != err {
		log.Fatal(err)
	}
	// window run
	window.Run(gm_env.GfxLoopTesting, probs)
}
