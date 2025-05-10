package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"us.figge.chess/internal/game"
)

func main() {
	g := game.NewGame()
	ebiten.SetWindowTitle("Lutefisk Chess Engine 2.0")
	err := ebiten.RunGame(g)
	if err != nil {
		log.Fatalf("Failed to initialize graphics engine: %v\n", err)
	}
	fmt.Println("Game: Done")
}
