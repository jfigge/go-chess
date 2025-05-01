package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"us.figge.chess/internal/game"
)

const (
	SheetImageSize = 426
	SquareSize     = SheetImageSize / 6
	ScreenWidth    = SquareSize * 8
	ScreenHeight   = SquareSize*8 + 20
)

func main() {
	g := game.NewGame(
		game.OptSquareSize(SquareSize),
		game.OptSheetImageSize(SheetImageSize),
		game.OptWhiteRGB(0xf1, 0xd9, 0xc0),
		game.OptBlackRGB(0xa9, 0x7a, 0x65),
	)
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Lutefisk Chess Engine 1.0")

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
