package board

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"strconv"
	"us.figge.chess/internal/board/colors"
	"us.figge.chess/internal/board/graphics"
	"us.figge.chess/internal/board/highligher"
	. "us.figge.chess/internal/common"
	"us.figge.chess/internal/engine"
)

const (
	debugHeight = 16
)

type Board struct {
	colors     *colors.Colors
	engine     *engine.Engine
	squareSize int

	// Graphics elements
	canvas        *ebiten.Image
	background    *ebiten.Image
	highlights    *ebiten.Image
	foreground    *ebiten.Image
	foregroundOp  [64]*ebiten.DrawImageOptions
	overlay       *ebiten.Image
	labelingX     *ebiten.Image
	labelingY     *ebiten.Image
	labelingXOp   *ebiten.DrawImageOptions
	labelFontSize float64

	// Highlighting
	redraw      bool
	regenerate  bool
	rehighlight bool
	selector    *highligher.DragAndDrop
	dragStart   *highligher.Highlight
	enPassant   *highligher.Highlight
	validMoves  []*highligher.ValidMove

	// Debugging
	debugEnabled bool
	lastCursorX  int
	lastCursorY  int
	debugY       int
	debugX       [8]int
}

func NewBoard(engine *engine.Engine, options ...Option) *Board {
	b := &Board{
		colors:     colors.NewColors(),
		engine:     engine,
		squareSize: 71,
	}
	for _, option := range options {
		option(b)
	}
	b.Initialize()
	return b
}

func (b *Board) Initialize() {
	graphics.InitPieces(b.squareSize)
	b.labelFontSize = float64(b.squareSize) * .17
	w, _ := graphics.TextSize("8", b.labelFontSize)
	_, h := graphics.TextSize("8", b.labelFontSize-2)
	b.canvas = ebiten.NewImage(b.squareSize*8, b.squareSize*8)
	b.background = ebiten.NewImage(b.squareSize*8, b.squareSize*8)
	b.highlights = ebiten.NewImage(b.squareSize*8, b.squareSize*8)
	b.foreground = ebiten.NewImage(b.squareSize*8, b.squareSize*8)
	b.overlay = ebiten.NewImage(b.squareSize*8, b.squareSize*8)
	b.labelingX = ebiten.NewImage(b.squareSize*8, int(h))
	b.labelingXOp = &ebiten.DrawImageOptions{}
	b.labelingXOp.GeoM.Translate(0, float64(b.squareSize*8)-h)
	b.labelingY = ebiten.NewImage(int(w), b.squareSize*8)
	b.selector = highligher.NewDragAndDrop(b, b.squareSize, b.colors.Tints(b.colors.Highlight()), b.colors.Tints(b.colors.Valid()))
	b.dragStart = highligher.NewHighlight(b, b.squareSize, b.colors.Tints(b.colors.DragStart()))
	b.enPassant = highligher.NewHighlight(b, b.squareSize, b.colors.Tints(b.colors.EnPassant()))

	for i := range 8 {
		b.debugY = b.squareSize*8 + 1
		b.debugX[i] = b.squareSize*i + 1
	}
	height := b.squareSize * 8
	if b.debugEnabled {
		height += debugHeight + 2
	}
	ebiten.SetWindowSize(b.squareSize*8, height)

	b.generateBackground()
}

func (b *Board) Update() error {
	x, y := ebiten.CursorPosition()
	b.lastCursorX, b.lastCursorY = x-1, y-2
	b.rehighlight = b.rehighlight || b.selector.Update(b.lastCursorX, b.lastCursorY)
	return nil
}

func (b *Board) Draw(screen *ebiten.Image) {
	if b.rehighlight || b.redraw {
		b.highlights.Clear()
		b.selector.Draw(b.highlights)
		b.enPassant.Draw(b.highlights)
		for i := range b.validMoves {
			b.validMoves[i].Draw(b.highlights)
		}
		b.rehighlight = false
		b.redraw = true
	}

	if b.redraw {
		b.canvas.DrawImage(b.background, nil)
		b.canvas.DrawImage(b.highlights, nil)
		b.canvas.DrawImage(b.labelingX, b.labelingXOp)
		b.canvas.DrawImage(b.labelingY, nil)
		b.canvas.DrawImage(b.foreground, nil)
		b.canvas.DrawImage(b.overlay, nil)
		b.dragStart.Draw(b.canvas)
		b.redraw = false
	}
	screen.DrawImage(b.canvas, nil)
	if b.selector.IsDragging() {
		b.selector.DrawDrag(screen)
	}
	if b.debugEnabled {
		s := float32(b.squareSize * 8)
		vector.DrawFilledRect(screen, 0, s, s, float32(b.debugY), b.colors.Black(), false)
		if b.selector != nil {
			b.selector.Debug(screen, b.debugX, b.debugY)
		}
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("X,Y:%d,%d", b.lastCursorX, b.lastCursorY), b.debugX[6], b.debugY)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("  TPS: %0.0f", ebiten.ActualTPS()), b.debugX[7], b.debugY)
	}
}

func (b *Board) Setup(fen string) {
	b.engine.SetFEN(fen)
	b.generateForeground()
}

func (b *Board) GetPieceType(rank, file uint8) (uint8, bool) {
	return b.engine.GetPieceType(rank, file)
}

func (b *Board) DragBegin(index, pieceType uint8) {
	rank, file := ItoRF(index)
	b.dragStart.Update(RFtoXY(rank, file, b.squareSize))
	b.updateValidMoves(index, pieceType)
	b.generateForeground()
}

func (b *Board) DragOver(index, pieceType uint8) {
	b.updateValidMoves(index, pieceType)
	b.rehighlight = true
}

func (b *Board) DragEnd(from, to, pieceType uint8, cancelled bool) {
	b.dragStart.Hide()
	b.rehighlight = true
	b.validMoves = nil
	if !cancelled {
		msg, ok := b.engine.MovePiece(from, to, pieceType)
		if ok {
			fmt.Printf("Move: %s\n", msg)
			if epi, ok := b.engine.GetEnPassant(); ok {
				b.enPassant.UpdateByIndex(epi)
			} else {
				b.enPassant.Hide()
			}
		}
	}
	b.generateForeground()
}

func (b *Board) updateValidMoves(index, pieceType uint8) {
	b.validMoves = nil
	rank, file := ItoRF(index)
	getMoves := b.engine.GetMoves(rank, file, pieceType)
	for i := range getMoves {
		valid := highligher.NewValidMove(b, b.squareSize, b.colors.Tints(b.colors.Valid()), getMoves[i])
		b.validMoves = append(b.validMoves, valid)
	}
}

func (b *Board) generateForeground() {
	dragIndex := -1
	if b.selector.IsDragging() {
		dragIndex = int(b.selector.DragIndex())
	}
	b.foreground.Clear()
	bitBoards := b.engine.GetBoards()
bit:
	for i := range 64 {
		if i == dragIndex {
			continue
		}
		bit := uint64(1 << (63 - i))
		player := PlayerWhite
		if bitBoards[BitBlack]&bit != 0 {
			player = PlayerBlack
		}
		var piece *graphics.Piece
		switch {
		case bitBoards[BitPawns]&bit != 0:
			piece = graphics.GetPiece(PiecePawn | player)
		case bitBoards[BitKnights]&bit != 0:
			piece = graphics.GetPiece(PieceKnight | player)
		case bitBoards[BitBishops]&bit != 0:
			piece = graphics.GetPiece(PieceBishop | player)
		case bitBoards[BitRooks]&bit != 0:
			piece = graphics.GetPiece(PieceRook | player)
		case bitBoards[BitQueens]&bit != 0:
			piece = graphics.GetPiece(PieceQueen | player)
		case bitBoards[BitKings]&bit != 0:
			piece = graphics.GetPiece(PieceKing | player)
		default:
			continue bit
		}
		piece.Draw(b.foreground, b.foregroundOp[i])
	}
	b.regenerate = false
	b.redraw = true
}

func (b *Board) generateBackground() {
	s := float32(b.squareSize)
	oddEven := 0
	clr := []color.Color{b.colors.PlayerWhite(), b.colors.PlayerBlack()}

	w, _ := graphics.TextSize("8", b.labelFontSize-2)
	for i := range 8 {
		for j := range 8 {
			index := i*8 + j
			op := ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(j*b.squareSize), float64(i*b.squareSize))
			b.foregroundOp[index] = &op
			vector.DrawFilledRect(b.background, float32(i)*s, float32(j)*s, s, s, clr[oddEven], false)
			oddEven = 1 - oddEven
		}
		wp, _ := graphics.TextSize(strconv.Itoa(i+1), b.labelFontSize-2)

		graphics.TextAt(b.labelingX, string([]byte{byte('A' + i)}), (i+1)*b.squareSize-int(w*1.5), 0, b.labelFontSize, clr[oddEven])
		graphics.TextAt(b.labelingY, strconv.Itoa(i+1), int((w-wp)/2), (7-i)*b.squareSize, b.labelFontSize, clr[oddEven])
		oddEven = 1 - oddEven
	}
}
