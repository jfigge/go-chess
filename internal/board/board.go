package board

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"strconv"
	"us.figge.chess/internal/board/graphics"
	. "us.figge.chess/internal/common"
	"us.figge.chess/internal/engine"
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

	colorWhite     color.Color
	colorBlack     color.Color
	colorValid     color.Color
	colorInvalid   color.Color
	colorHighlight color.Color

	pieces *graphics.Pieces

	squareSize int
	fontHeight int
}

func NewBoard(engine *engine.Engine, options ...Option) *Board {
	b := &Board{
		engine:         engine,
		squareSize:     71,
		fontHeight:     16,
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
	b.labelFontSize = float64(b.squareSize) * .17
	w, _ := graphics.TextSize("8", b.labelFontSize)
	_, h := graphics.TextSize("8", b.labelFontSize-2)
	b.pieces = graphics.NewPieces(b.squareSize)
	b.canvas = ebiten.NewImage(b.squareSize*8, b.squareSize*8)
	b.background = ebiten.NewImage(b.squareSize*8, b.squareSize*8)
	b.highlights = ebiten.NewImage(b.squareSize*8, b.squareSize*8)
	b.foreground = ebiten.NewImage(b.squareSize*8, b.squareSize*8)
	b.overlay = ebiten.NewImage(b.squareSize*8, b.squareSize*8)
	b.labelingX = ebiten.NewImage(b.squareSize*8, int(h))
	b.labelingXOp = &ebiten.DrawImageOptions{}
	b.labelingXOp.GeoM.Translate(0, float64(b.squareSize*8)-h)
	b.labelingY = ebiten.NewImage(int(w), b.squareSize*8)

	ebiten.SetWindowSize(b.squareSize*8, b.squareSize*8)
	b.generateBackground()
}

func (b *Board) Draw(screen *ebiten.Image) {
	b.canvas.DrawImage(b.background, nil)
	b.canvas.DrawImage(b.highlights, nil)
	b.canvas.DrawImage(b.labelingX, b.labelingXOp)
	b.canvas.DrawImage(b.labelingY, nil)
	b.canvas.DrawImage(b.foreground, nil)
	b.canvas.DrawImage(b.overlay, nil)
	screen.DrawImage(b.canvas, nil)

}

func (b *Board) Setup(fen string) {
	b.engine.Setup(fen)
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
			piece = b.pieces.Piece(PiecePawn | player)
		case bitBoards[engine.BitKnights]&bit != 0:
			piece = b.pieces.Piece(PieceKnight | player)
		case bitBoards[engine.BitBishops]&bit != 0:
			piece = b.pieces.Piece(PieceBishop | player)
		case bitBoards[engine.BitRooks]&bit != 0:
			piece = b.pieces.Piece(PieceRook | player)
		case bitBoards[engine.BitQueens]&bit != 0:
			piece = b.pieces.Piece(PieceQueen | player)
		case bitBoards[engine.BitKings]&bit != 0:
			piece = b.pieces.Piece(PieceKing | player)
		default:
			continue bit
		}
		piece.Draw(b.foreground, b.foregroundOp[i])
	}
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
