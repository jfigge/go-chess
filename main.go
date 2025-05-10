package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"os"
	"us.figge.chess/internal/game"
)

var (
	GitCommit string
	GitBranch string
	GitTag    string
	Built     string
	GitSource string
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("Lutefisk Chess Engine 2.0 - by Jason Figge\n")
		fmt.Printf("  Git source: %s\n", GitSource)
		fmt.Printf("  Git commit: %s\n", GitCommit)
		fmt.Printf("  Git branch: %s\n", GitBranch)
		fmt.Printf("  Git tag:    %s\n", GitTag)
		fmt.Printf("  Built:      %s\n\n", Built)
		return
	}
	g := game.NewGame()
	ebiten.SetWindowTitle("Lutefisk Chess Engine 2.0")
	err := ebiten.RunGame(g)
	if err != nil {
		log.Fatalf("Failed to initialize graphics engine: %v\n", err)
	}
	fmt.Println("Game: Done")
}
