package board

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"strings"
	. "us.figge.chess/internal/common"
)

type highlight struct {
	Configuration
	cursorX  int
	cursorY  int
	index    int
	rank     int
	file     int
	notation string
	x        float64
	y        float64
	size     float32
	piece    Piece
	color    color.Color
}

type Board struct {
	Configuration
	position
	ui
	fen       string
	highlight *highlight

	//dragPiece *piece.Piece
	//dragIndex int
}

func NewBoard(c Configuration, options ...Options) *Board {
	board := &Board{
		Configuration: c,
		position:      position{Configuration: c},
	}
	for _, option := range options {
		option(board)
	}
	board.initializeImages()
	board.renderBackground()
	board.setupBoard(board.fen)
	fmt.Println("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	fmt.Println(board.generateFen())
	board.renderForeground()
	return board
}

func (b *Board) Update() {
	x, y := ebiten.CursorPosition()
	rank, file, ok := b.TranslateXYtoRF(x, y)
	if ok {
		index := b.TranslateRFtoIndex(rank, file)
		b.highlight = &highlight{
			Configuration: b.Configuration,
			cursorX:       x,
			cursorY:       y,
			rank:          rank,
			file:          file,
			index:         index,
			notation:      b.TranslateRFtoN(rank, file),
			size:          float32(b.SquareSize()),
			color:         b.ColorHighlight(),
		}
		if pieceType, ok := b.findPiece(index); ok {
			b.highlight.piece = b.Piece(pieceType)
		}
		b.highlight.x, b.highlight.y = b.TranslateRFtoXY(rank, file)
	} else {
		b.highlight = nil
	}
	//if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && b.highlightSquare != nil && b.highlightSquare.piece != nil {
	//	b.dragIndex = b.highlightSquare.index
	//	b.dragPiece = b.highlightSquare.piece
	//	b.highlightSquare.piece = nil
	//	b.dragPiece.StartDrag(x, y)
	//	b.renderForeground()
	//	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	//} else if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && b.dragPiece != nil {
	//	ebiten.SetCursorMode(ebiten.CursorModeVisible)
	//	if ok && b.highlightSquare.piece == nil {
	//		b.dragPiece.StopDrag(rank, file)
	//		b.highlightSquare.piece = b.dragPiece
	//	} else {
	//		b.dragPiece.CancelDrag()
	//		//b.squares[b.dragIndex].piece = b.dragPiece
	//	}
	//	b.dragPiece = nil
	//	b.dragIndex = -1
	//	b.renderForeground()
	//}
	//if b.HighlightAttacks() {
	//
	//}
}

func (b *Board) Draw(target *ebiten.Image) {
	b.composite.DrawImage(b.background, nil)
	if b.highlight != nil {
		b.highlight.Draw(b.composite)
	}
	if b.ShowLabels() {
		b.composite.DrawImage(b.labelingX, b.labelingXOp)
		b.composite.DrawImage(b.labelingY, nil)
	}
	b.composite.DrawImage(b.foreground, nil)
	//if b.dragPiece != nil && b.highlightSquare != nil {
	//	b.dragPiece.Draw(b.composite, false)
	//}
	target.DrawImage(b.composite, nil)
	if b.ShowStrength() {
		target.DrawImage(b.strength, b.strengthOp)
	}
	//if b.enpassant != -1 {
	//	rank, file := b.TranslateIndexToRF(b.enpassant)
	//	x, y := b.TranslateRFtoXY(rank, file)
	//	vector.DrawFilledRect(target, float32(x), float32(y), float32(b.SquareSize()), float32(b.SquareSize()), b.ColorValid(), false)
	//}
	if b.EnableDebug() {
		if b.highlight != nil {
			b.highlight.Debug(target)
		}
		enPassant := b.EnPassant()
		if enPassant != 0 {
			rank := "6"
			if b.Turn() == Black {
				rank = "3"
			}
			ebitenutil.DebugPrintAt(target, fmt.Sprintf("EnPas: "+string([]byte{'a' + 8 - enPassant})+rank), b.DebugX(2), b.DebugY())
		}
		ebitenutil.DebugPrintAt(target, "Move: "+b.TurnName(b.Turn()), b.DebugX(6), b.DebugY())
	}
}
func (h *highlight) Draw(dst *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Reset()
	op.GeoM.Translate(h.x, h.y)
	vector.DrawFilledRect(dst, float32(h.x), float32(h.y), h.size, h.size, h.color, false)
}

func (h *highlight) Debug(dst *ebiten.Image) {
	ebitenutil.DebugPrintAt(dst, fmt.Sprintf("Index: %d", h.index), h.DebugX(0), h.DebugY())
	ebitenutil.DebugPrintAt(dst, fmt.Sprintf("RF: %d,%d %s", h.rank, h.file, strings.ToUpper(h.notation)), h.DebugX(1), h.DebugY())
	ebitenutil.DebugPrintAt(dst, fmt.Sprintf("XY: %d,%d", h.cursorX, h.cursorY), h.DebugX(3), h.DebugY())
	if h.piece != nil {
		ebitenutil.DebugPrintAt(dst, h.piece.Color()+" "+h.piece.Name(), h.DebugX(4), h.DebugY())
	}
}
