package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"us.figge.chess/internal/board"
)

type Game struct {
	*ColorScheme
	entities       map[uint8]*Entity
	board          *board.Board
	squareSize     uint
	sheetImageSize int
	enabledDebug   bool
	debugY         int
	debugX         [8]int
}

func NewGame(options ...GameOptions) *Game {
	game := &Game{
		ColorScheme:  newColorScheme(),
		squareSize:   64,
		enabledDebug: true,
	}
	for _, option := range options {
		option(game)
	}
	game.entities = makeEntities(game)
	game.board = board.NewBoard(
		game,
		//board.OptSetup("r1bk3r/p2pBpNp/n4n2/1p1NP2P/6P1/3P4/P1P1K3/q5b1"),
	)
	for i := 0; i < 8; i++ {
		game.debugX[i] = int(game.squareSize)*i + 2
	}
	game.debugY = int(game.squareSize*8 + 2)
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
	if g.EnableDebug() {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS()), g.debugX[7], g.debugY)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
