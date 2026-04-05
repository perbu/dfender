package main

import (
	_ "embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/perbu/dfender/game"
)

//go:embed assets/One_Last_Life_Remaining.mp3
var musicData []byte

func main() {
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("dFender")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g := game.New(musicData)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
