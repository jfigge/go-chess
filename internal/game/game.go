package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"us.figge.chess/internal/board"
	"us.figge.chess/internal/engine"
)

const (
	FontHeight = 16
	SquareSize = 71
)

type Game struct {
	board  *board.Board
	engine *engine.Engine
}

func NewGame() *Game {
	eng := engine.NewEngine()
	g := &Game{
		board: board.NewBoard(eng,
			board.OptSquareSize(SquareSize),
			board.OptFontHeight(FontHeight),
			board.OptWhiteRGB(0xf1, 0xd9, 0xc0),
			board.OptBlackRGB(0xa9, 0x7a, 0x65),
		),
	}

	g.board.Setup("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq _ 0 1")
	return g
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.board.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
