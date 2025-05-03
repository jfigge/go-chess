package board

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
	"strconv"
	"strings"
	. "us.figge.chess/internal/common"
	"us.figge.chess/internal/piece"
	"us.figge.chess/internal/player"
)

type Board struct {
	Configuration
	players [2]*player.Player

	squares          [64]*square
	composite        *ebiten.Image
	background       *ebiten.Image
	labelingX        *ebiten.Image
	labelingY        *ebiten.Image
	labelingXOp      *ebiten.DrawImageOptions
	labelingFontSize float64
	foreground       *ebiten.Image
	strength         *ebiten.Image
	strengthOp       *ebiten.DrawImageOptions
	highlightSquare  *square

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
		Configuration:    c,
		labelingFontSize: float64(c.SquareSize()) * .17,
		fen:              "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
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

func (b *Board) initializeImages() {
	w, _ := b.TextSize("8", b.labelingFontSize)
	_, h := b.TextSize("8", b.labelingFontSize-2)
	s := b.SquareSize()
	b.composite = ebiten.NewImage(s*8, s*8)
	b.background = ebiten.NewImage(s*8, s*8)
	b.labelingX = ebiten.NewImage(s*8, int(h))
	b.labelingXOp = &ebiten.DrawImageOptions{}
	b.labelingXOp.GeoM.Translate(0, float64(s*8)-h)
	b.labelingY = ebiten.NewImage(int(w), s*8)
	b.strength = ebiten.NewImage(b.FontHeight(), s*8)
	b.strengthOp = &ebiten.DrawImageOptions{}
	b.strengthOp.GeoM.Translate(float64(s*8), 0)
}

func (b *Board) renderBackground() {
	s := float32(b.SquareSize())
	oddEven := 0
	clr := []color.Color{b.ColorWhite(), b.ColorBlack()}

	_, h := b.TextSize("8", b.labelingFontSize-2)
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			vector.DrawFilledRect(b.background, float32(i)*s, float32(j)*s, s, s, clr[oddEven], false)
			oddEven = 1 - oddEven
		}
		if i == 0 {
			b.TextAt(b.labelingX, "A1", 0, 0, b.labelingFontSize, clr[oddEven])
		} else {
			b.TextAt(b.labelingX, string([]byte{byte('A' + i)}), i*b.SquareSize(), 0, b.labelingFontSize, clr[oddEven])
			b.TextAt(b.labelingY, strconv.Itoa(i+1), 0, (8-i)*b.SquareSize()-int(h), b.labelingFontSize, clr[oddEven])
		}
		oddEven = 1 - oddEven
	}
}

func (b *Board) renderForeground() {
	b.foreground = ebiten.NewImage(b.SquareSize()*8, b.SquareSize()*8)
	b.players[0].Draw(b.foreground)
	b.players[1].Draw(b.foreground)

	s := float32(b.SquareSize() * 8)
	_, y := ebiten.CursorPosition()
	pct := float32(y) / s
	h1 := s * pct
	val := int(math.Abs(float64(pct-.5) * 200))
	if val > 100 {
		val = 100
	}
	b.strength.Clear()
	vector.DrawFilledRect(b.strength, 0, h1, float32(b.FontHeight()), s-h1, b.ColorStrength(), false)
	ebitenutil.DebugPrintAt(b.strength, fmt.Sprintf("%d", val), 0, int((s-float32(b.FontHeight()))/2))

	b.fen = b.Fen()
}
