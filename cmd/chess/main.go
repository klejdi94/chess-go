package main

import (
	"github.com/user/chess/pkg/game"
	"github.com/user/chess/pkg/ui"
)

func main() {
	game := game.NewGame()
	ui := ui.NewUI(game)
	ui.Start()
}
