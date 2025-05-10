package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"us.figge.chess/internal/board"
	"us.figge.chess/internal/engine"
)

const (
	SquareSize = 71
)

type Game struct {
	board  *board.Board
	engine *engine.Engine
}

func NewGame() *Game {
	ebiten.SetTPS(20)
	ebiten.SetVsyncEnabled(true)
	ebiten.SetScreenClearedEveryFrame(false)
	eng := engine.NewEngine()
	g := &Game{
		board: board.NewBoard(eng,
			board.OptSquareSize(SquareSize),
			board.OptDebugEnabled(true),
			board.OptWhiteRGB(0xf1, 0xd9, 0xc0),
			board.OptBlackRGB(0xa9, 0x7a, 0x65),
		),
	}
	g.board.Setup("")
	return g
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}
	return g.board.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.board.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
