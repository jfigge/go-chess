package board

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"strconv"
	"us.figge.chess/internal/board/graphics"
	. "us.figge.chess/internal/common"
	"us.figge.chess/internal/engine"
)

const (
	debugHeight = 16
)

type Board struct {
	engine        *engine.Engine
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

	redraw      bool
	rehighlight bool
	selector    *Highlight

	debugEnabled bool
	debugY       int
	debugX       [8]int

	colorWhite     color.Color
	colorBlack     color.Color
	colorValid     color.Color
	colorInvalid   color.Color
	colorHighlight color.Color

	squareSize int
}

func NewBoard(engine *engine.Engine, options ...Option) *Board {
	b := &Board{
		engine:         engine,
		squareSize:     71,
		colorWhite:     &color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		colorBlack:     &color.RGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xff},
		colorValid:     &color.RGBA{R: 0x44, G: 0xff, B: 0x44, A: 0xff},
		colorInvalid:   &color.RGBA{R: 0xff, G: 0x44, B: 0x44, A: 0xff},
		colorHighlight: &color.RGBA{R: 0x88, G: 0x22, B: 0x22, A: 0xff},
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
	b.selector = NewHighlight(b.engine, b.squareSize, b.colorHighlight)

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
	b.rehighlight = b.rehighlight || b.selector.Update(x, y)
	return nil
}

func (b *Board) Draw(screen *ebiten.Image) {
	if b.rehighlight {
		b.highlights.Clear()
		b.selector.Draw(b.highlights)
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
		b.redraw = false
	}
	screen.DrawImage(b.canvas, nil)
	if b.debugEnabled {
		if b.selector != nil {
			b.selector.Debug(screen, b.debugX, b.debugY)
		}
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS()), b.debugX[7], b.debugY)
	}
}

func (b *Board) Setup(fen string) {
	b.engine.SetFEN(fen)
	b.generateForeground()
}

func (b *Board) generateForeground() {
	b.foreground.Clear()
	bitBoards := b.engine.GetBoards()
bit:
	for i := range 64 {
		bit := uint64(1 << i)
		player := PlayerWhite
		if bitBoards[engine.BitBlack]&bit != 0 {
			player = PlayerBlack
		}
		var piece *graphics.Piece
		switch {
		case bitBoards[engine.BitPawns]&bit != 0:
			piece = graphics.GetPiece(PiecePawn | player)
		case bitBoards[engine.BitKnights]&bit != 0:
			piece = graphics.GetPiece(PieceKnight | player)
		case bitBoards[engine.BitBishops]&bit != 0:
			piece = graphics.GetPiece(PieceBishop | player)
		case bitBoards[engine.BitRooks]&bit != 0:
			piece = graphics.GetPiece(PieceRook | player)
		case bitBoards[engine.BitQueens]&bit != 0:
			piece = graphics.GetPiece(PieceQueen | player)
		case bitBoards[engine.BitKings]&bit != 0:
			piece = graphics.GetPiece(PieceKing | player)
		default:
			continue bit
		}
		piece.Draw(b.foreground, b.foregroundOp[i])
	}
	b.redraw = true
}

func (b *Board) generateBackground() {
	s := float32(b.squareSize)
	oddEven := 0
	clr := []color.Color{b.colorWhite, b.colorBlack}

	_, h := graphics.TextSize("8", b.labelFontSize-2)
	for i := range 8 {
		for j := range 8 {
			index := i*8 + j
			op := ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(j*b.squareSize), float64((7-i)*b.squareSize))
			b.foregroundOp[index] = &op
			vector.DrawFilledRect(b.background, float32(i)*s, float32(j)*s, s, s, clr[oddEven], false)
			oddEven = 1 - oddEven
		}
		if i == 0 {
			graphics.TextAt(b.labelingX, "A1", 0, 0, b.labelFontSize, clr[oddEven])
		} else {
			graphics.TextAt(b.labelingX, string([]byte{byte('A' + i)}), i*b.squareSize, 0, b.labelFontSize, clr[oddEven])
			graphics.TextAt(b.labelingY, strconv.Itoa(i+1), 0, (8-i)*b.squareSize-int(h), b.labelFontSize, clr[oddEven])
		}
		oddEven = 1 - oddEven
	}
}
