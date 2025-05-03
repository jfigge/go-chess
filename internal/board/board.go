package board

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"strconv"
	"us.figge.chess/internal/player"
	. "us.figge.chess/internal/shared"
)

type Board struct {
	Configuration
	players    [2]*player.Player
	squares    [64]*square
	background *ebiten.Image
	foreground *ebiten.Image

	highlightSquare *square

	turn      uint8
	enpassant int
	fullMove  int
	halfMove  int
	fen       string

	//dragPiece      *piece.Piece
	//dragIndex      int
}

func NewBoard(c Configuration, options ...BoardOptions) *Board {
	board := &Board{
		Configuration: c,
		background:    renderBoardBackground(c),
		fen:           "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
	}
	for _, option := range options {
		option(board)
	}
	board.SetFen(board.fen)
	return board
}

func (b *Board) Update() {
	x, y := ebiten.CursorPosition()
	if rank, file, ok := b.TranslateXYtoRF(x, y); ok {
		b.highlightSquare = b.squares[b.TranslateRFtoIndex(rank, file)]
	} else {
		b.highlightSquare = nil
	}
}

func (b *Board) Draw(target *ebiten.Image) {
	target.DrawImage(b.background, nil)
	if b.highlightSquare != nil {
		b.highlightSquare.Draw(b, target)
	}
	target.DrawImage(b.foreground, nil)

	if b.EnableDebug() {
		x, y := ebiten.CursorPosition()
		ebitenutil.DebugPrintAt(target, fmt.Sprintf("X,Y: %d,%d", x, y), b.DebugX(2), b.DebugY())
		if r, f, ok := b.TranslateXYtoRF(x, y); ok {
			ebitenutil.DebugPrintAt(target, fmt.Sprintf("R,F: %d,%d", r, f), b.DebugX(1), b.DebugY())
		}
		if b.highlightSquare != nil {
			index := b.highlightSquare.index
			ebitenutil.DebugPrintAt(target, "Index: "+strconv.Itoa(index), b.DebugX(0), b.DebugY())
			p := b.highlightSquare.piece
			if p != nil {
				ebitenutil.DebugPrintAt(target, p.Token.Color()+" "+p.Name(), b.DebugX(4), b.DebugY())
			}
		}
		ebitenutil.DebugPrintAt(target, "Move: "+b.Color(b.turn), b.DebugX(6), b.DebugY())
		ebitenutil.DebugPrintAt(target, "Fen:"+b.fen, b.DebugX(0), b.DebugFen())
	}
}

func renderBoardBackground(c Configuration) *ebiten.Image {
	s := c.SquareSize()
	k := 0
	clr := []color.Color{c.ColorWhite(), c.ColorBlack()}
	img := ebiten.NewImage(s*8, s*8)
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			vector.DrawFilledRect(img, float32(i*s), float32(j*s), float32(s), float32(s), clr[k], false)
			k = 1 - k
		}
		k = 1 - k
	}
	return img
}
