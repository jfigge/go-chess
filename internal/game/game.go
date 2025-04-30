package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"us.figge.chess/internal/board"
	"us.figge.chess/internal/shared"
)

type GameOptions func(g *Game)

func OptSquareSize(size uint) GameOptions {
	return func(g *Game) {
		g.squareSize = size
	}
}
func OptSheetImageSize(size int) GameOptions {
	return func(g *Game) {
		g.sheetImageSize = size
	}
}

type Game struct {
	*ColorScheme
	entities       map[uint8]*Entity
	board          *board.Board
	squareSize     uint
	sheetImageSize int
}

func NewGame(options ...GameOptions) *Game {
	game := &Game{
		ColorScheme: newColorScheme(),
		squareSize:  64,
	}
	for _, option := range options {
		option(game)
	}
	game.entities = makeEntities(game)
	game.board = board.NewBoard(
		game,
		//board.OptSetup("r1bk3r/p2pBpNp/n4n2/1p1NP2P/6P1/3P4/P1P1K3/q5b1"),
	)
	return game
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		_ = ebiten.Termination
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.board.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *Game) SquareSize() uint                   { return g.squareSize }
func (g *Game) Token(pieceType uint8) shared.Token { return g.entities[pieceType] }
func (g *Game) SheetImageSize() int                { return g.sheetImageSize }
func (g *Game) Translate(rank, file uint8) (float64, float64) {
	return float64(uint(rank-1) * g.squareSize), float64(uint(8-file) * g.squareSize)
}
