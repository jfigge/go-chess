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
	squareSize     int
	sheetImageSize int

	fontHeight   int
	boardWidth   int
	boardHeight  int
	debugEnabled bool
	debugY       int
	debugX       [8]int
	fenY         int
}

func NewGame(options ...GameOptions) *Game {
	game := &Game{
		ColorScheme:  newColorScheme(),
		squareSize:   64,
		boardWidth:   512,
		boardHeight:  512,
		debugEnabled: false,
		fontHeight:   16,
	}
	for _, option := range options {
		option(game)
	}
	game.entities = makeEntities(game)
	game.board = board.NewBoard(
		game,
		board.OptSetup("rn1qkbnr/pppp1ppp/b2Qp3/8/8/8/PPPPPPPP/RNB1KBNR b Qq c4"),
		//board.OptSetup("r1bk3r/p2pBpNp/n4n2/1p1NP2P/6P1/3P4/P1P1K3/q5b1"),
	)
	for i := 0; i < 8; i++ {
		game.debugX[i] = game.squareSize*i + 2
	}
	game.debugY = game.squareSize*8 + 2
	game.fenY = game.squareSize*8 + game.fontHeight + 2
	game.boardWidth = game.squareSize * 8
	game.boardHeight = game.squareSize * 8
	if game.debugEnabled {
		game.boardHeight += game.fontHeight*2 + 2
	}
	ebiten.SetWindowSize(game.boardWidth, game.boardHeight)

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
