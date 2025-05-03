package board

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"strings"
	. "us.figge.chess/internal/common"
	"us.figge.chess/internal/piece"
	"us.figge.chess/internal/player"
)

type Board struct {
	Configuration
	ui
	players   [2]*player.Player
	squares   [64]*square
	turn      uint8
	castling  uint8
	enpassant int
	fullMove  int
	halfMove  int
	fen       string

	dragPiece *piece.Piece
	dragIndex int
}

func NewBoard(c Configuration, options ...Options) *Board {
	board := &Board{
		Configuration: c,
	}
	for _, option := range options {
		option(board)
	}
	board.initializeImages()
	board.renderBackground()
	board.SetFen(board.fen)
	return board
}

func (b *Board) Update() {
	x, y := ebiten.CursorPosition()
	rank, file, ok := b.TranslateXYtoRF(x, y)
	if ok {
		b.highlightSquare = b.squares[b.TranslateRFtoIndex(rank, file)]
	} else {
		b.highlightSquare = nil
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && b.highlightSquare != nil && b.highlightSquare.piece != nil {
		b.dragIndex = b.highlightSquare.index
		b.dragPiece = b.highlightSquare.piece
		b.highlightSquare.piece = nil
		b.dragPiece.StartDrag(x, y)
		b.renderForeground()
		ebiten.SetCursorMode(ebiten.CursorModeHidden)
	} else if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && b.dragPiece != nil {
		ebiten.SetCursorMode(ebiten.CursorModeVisible)
		if ok && b.highlightSquare.piece == nil {
			b.dragPiece.StopDrag(rank, file)
			b.highlightSquare.piece = b.dragPiece
		} else {
			b.dragPiece.CancelDrag()
			b.squares[b.dragIndex].piece = b.dragPiece
		}
		b.dragPiece = nil
		b.dragIndex = -1
		b.renderForeground()
	}
	if b.HighlightAttacks() {

	}
}

func (b *Board) Draw(target *ebiten.Image) {
	b.composite.DrawImage(b.background, nil)
	if b.highlightSquare != nil {
		b.highlightSquare.Draw(b, b.composite)
	}
	if b.ShowLabels() {
		b.composite.DrawImage(b.labelingX, b.labelingXOp)
		b.composite.DrawImage(b.labelingY, nil)
	}
	b.composite.DrawImage(b.foreground, nil)
	if b.dragPiece != nil && b.highlightSquare != nil {
		b.dragPiece.Draw(b.composite, false)
	}
	target.DrawImage(b.composite, nil)
	if b.ShowStrength() {
		target.DrawImage(b.strength, b.strengthOp)
	}
	if b.enpassant != -1 {
		rank, file := b.TranslateIndexToRF(b.enpassant)
		x, y := b.TranslateRFtoXY(rank, file)
		vector.DrawFilledRect(target, float32(x), float32(y), float32(b.SquareSize()), float32(b.SquareSize()), b.ColorValid(), false)
	}
	if b.EnableDebug() {
		x, y := ebiten.CursorPosition()
		if r, f, ok := b.TranslateXYtoRF(x, y); ok {
			index := b.TranslateRFtoIndex(r, f)
			notation := b.TranslateRFtoN(r, f)
			ebitenutil.DebugPrintAt(target, fmt.Sprintf("Index: %d", index), b.DebugX(0), b.DebugY())
			ebitenutil.DebugPrintAt(target, fmt.Sprintf("RF: %d,%d %s", r, f, strings.ToUpper(notation)), b.DebugX(1), b.DebugY())
			if b.enpassant != -1 {
				rank, file := b.TranslateIndexToRF(b.enpassant)
				notation = b.TranslateRFtoN(rank, file)
				ebitenutil.DebugPrintAt(target, fmt.Sprintf("EnPas: "+notation), b.DebugX(2), b.DebugY())
			}
			ebitenutil.DebugPrintAt(target, fmt.Sprintf("XY: %d,%d", x, y), b.DebugX(3), b.DebugY())
			p := b.highlightSquare.piece
			if p != nil {
				ebitenutil.DebugPrintAt(target, p.Token.Color()+" "+p.Name(), b.DebugX(4), b.DebugY())
			}
		}
		ebitenutil.DebugPrintAt(target, "Move: "+b.Turn(b.turn), b.DebugX(6), b.DebugY())
	}
}

func (b *Board) Setup(fen string) {
	b.SetFen(fen)
}
