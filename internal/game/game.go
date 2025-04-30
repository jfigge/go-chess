package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"us.figge.chess/internal/board"
	"us.figge.chess/internal/shared"
)

type GameOptions func(g *Game)

func OptSquareSize(size int) GameOptions {
	return func(g *Game) {
		g.squareSize = size
	}
}
func OptSheetImageSize(size int) GameOptions {
	return func(g *Game) {
		g.sheetImageSize = size
	}
}
func OptWhiteRGB(red, green, blue uint8) GameOptions {
	return func(g *Game) {
		g.colorWhite = &color.RGBA{R: red, G: green, B: blue, A: 0xff}
	}
}
func OptBlackRGB(red, green, blue uint8) GameOptions {
	return func(g *Game) {
		g.colorBlack = &color.RGBA{R: red, G: green, B: blue, A: 0xff}
	}
}

type Game struct {
	entities       map[uint]*Entity
	board          *board.Board
	squareSize     int
	colorWhite     color.Color
	colorBlack     color.Color
	sheetImageSize int
}

func NewGame(options ...GameOptions) *Game {
	game := &Game{
		squareSize: 64,
		colorWhite: &color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		colorBlack: &color.RGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xff},
	}
	for _, option := range options {
		option(game)
	}
	game.entities = makeEntities(game)
	game.board = board.NewBoard(game)
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

func (g *Game) SquareSize() int                   { return g.squareSize }
func (g *Game) WhiteColor() color.Color           { return g.colorWhite }
func (g *Game) BlackColor() color.Color           { return g.colorBlack }
func (g *Game) Token(pieceType uint) shared.Token { return g.entities[pieceType] }
func (g *Game) SheetImageSize() int               { return g.sheetImageSize }
func (g *Game) Translate(rank, file int) (float64, float64) {
	return float64((file - 1) * g.squareSize), float64((rank - 1) * g.squareSize)
}
