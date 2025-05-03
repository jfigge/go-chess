package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"us.figge.chess/internal/board"
)

type Game struct {
	*ColorScheme
	entities       map[uint8]*Entity
	board          *board.Board
	squareSize     int
	sheetImageSize int

	highlightAttacks bool
	showStrength     bool
	showFPS          bool
	showLabels       bool
	fontHeight       int
	boardWidth       int
	boardHeight      int
	targetWidth      int
	targetHeight     int
	debugEnabled     bool
	debugY           int
	debugX           [8]int
}

func NewGame(options ...Options) *Game {
	game := &Game{
		ColorScheme:      newColorScheme(),
		squareSize:       64,
		boardWidth:       512,
		boardHeight:      512,
		targetWidth:      512,
		targetHeight:     512,
		debugEnabled:     false,
		highlightAttacks: false,
		showStrength:     false,
		showFPS:          false,
		fontHeight:       16,
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
	game.debugY = game.squareSize * 8
	game.boardWidth = game.squareSize * 8
	game.boardHeight = game.squareSize * 8
	if game.debugEnabled {
		game.boardHeight += game.fontHeight
	}
	if game.showStrength {
		game.boardWidth += game.fontHeight
	}
	game.targetWidth = game.boardWidth
	game.targetHeight = game.boardHeight
	ebiten.SetWindowSize(game.boardWidth, game.boardHeight)

	return game
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		return ebiten.Termination
	} else if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.debugEnabled = !g.debugEnabled
		if g.debugEnabled {
			g.targetHeight = g.squareSize*8 + g.fontHeight
		} else {
			g.targetHeight = g.squareSize * 8
		}
		ebiten.SetWindowSize(g.boardWidth, g.boardHeight)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.showStrength = !g.showStrength
		if g.showStrength {
			g.targetWidth = g.squareSize*8 + g.fontHeight
		} else {
			g.targetWidth = g.squareSize * 8
		}
		ebiten.SetWindowSize(g.boardWidth, g.boardHeight)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.highlightAttacks = !g.highlightAttacks
	} else if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		g.showFPS = !g.showFPS
	} else if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		g.showLabels = !g.showLabels
	}
	if g.targetHeight > g.boardHeight {
		g.boardHeight++
		ebiten.SetWindowSize(g.boardWidth, g.boardHeight)
	} else if g.targetHeight < g.boardHeight {
		g.boardHeight--
		ebiten.SetWindowSize(g.boardWidth, g.boardHeight)
	}
	if g.targetWidth > g.boardWidth {
		g.boardWidth++
		ebiten.SetWindowSize(g.boardWidth, g.boardHeight)
	} else if g.targetWidth < g.boardWidth {
		g.boardWidth--
		ebiten.SetWindowSize(g.boardWidth, g.boardHeight)
	}
	g.board.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.board.Draw(screen)
	if g.EnableDebug() && g.showFPS {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS()), g.DebugX(7), g.debugY)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
