package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	highlight      *ebiten.Image
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
	game.highlight = ebiten.NewImage(int(game.SquareSize()), int(game.SquareSize()))
	game.highlight.Fill(game.ColorHighlight())
	return game
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}
	g.board.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.board.Draw(screen)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS()), 0, 0)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *Game) SquareSize() uint                   { return g.squareSize }
func (g *Game) Token(pieceType uint8) shared.Token { return g.entities[pieceType] }
func (g *Game) SheetImageSize() int                { return g.sheetImageSize }
func (g *Game) TranslateRFtoXY(rank, file uint8) (float64, float64) {
	return float64(uint(rank-1) * g.squareSize), float64(uint(8-file) * g.squareSize)
}
func (g *Game) TranslateXYtoRF(x, y int) (uint8, uint8, bool) {
	rank := uint8(float32(x-1)/float32(g.squareSize)) + 1
	file := 8 - uint8(float32(y-2)/float32(g.squareSize))
	if rank < 1 || rank > 8 || file < 1 || file > 8 {
		return 0, 0, false
	}
	return rank, file, true
}
func (g *Game) TranslateRFtoIndex(rank, file uint8) uint8 {
	index := (8-file)*8 + rank - 1
	return index
}
func (g *Game) TranslateIndexToRF(index uint8) (uint8, uint8) {
	rank := index%8 + 1
	file := 8 - index/8
	return rank, file
}
func (g *Game) TranslateIndexToXY(index uint8) (float64, float64) {
	rank, file := g.TranslateIndexToRF(index)
	x, y := g.TranslateRFtoXY(rank, file)
	return x, y
}
