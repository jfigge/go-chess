package highlighers

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
	"us.figge.chess/internal/board/graphics"
	. "us.figge.chess/internal/common"
)

var (
	// whiteSubImage is an internal sub image of whiteImage.
	// Use whiteSubImage at DrawTriangles instead of whiteImage in order to avoid bleeding edges.
	whiteImage    = ebiten.NewImage(3, 3)
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	whiteImage.Fill(color.White)
}

type EnPassant struct {
	*Highlight
	path  *vector.Path
	piece *graphics.Piece
	vs    []ebiten.Vertex
	is    []uint16
	op    *ebiten.DrawTrianglesOptions
	x0    int
	x1    int
	x2    int
	y0    int
}

func NewEnPassant(highlighter Highlighter, squareSize int, validColor [2]color.Color, player uint8) *EnPassant {
	ep := &EnPassant{
		Highlight: NewHighlight(highlighter, squareSize, validColor),
	}
	ep.op = &ebiten.DrawTrianglesOptions{}
	ep.op.ColorScaleMode = ebiten.ColorScaleModePremultipliedAlpha
	ep.op.FillRule = ebiten.FillRuleFillAll
	ep.op.AntiAlias = true
	ep.x0 = squareSize / 2
	ep.x1 = squareSize / 4
	ep.x2 = ep.x1 / 5
	ep.y0 = squareSize / 5
	ep.piece = graphics.GetPiece(PiecePawn + player)
	return ep
}

func (h *EnPassant) UpdateByIndex(index uint8) {
	fmt.Println(index)
	h.y0 = h.squareSize / 5
	h.Highlight.UpdateByIndex(index)
	if index < 32 {
		h.y0 = -h.y0
		//h.highlightY += h.squareSize
	}
	h.path = &vector.Path{}
	h.path.MoveTo(float32(h.highlightX+h.x0), float32(h.highlightY+h.y0))
	h.path.LineTo(float32(h.highlightX+h.x0-h.x1), float32(h.highlightY+h.y0*3))
	h.path.LineTo(float32(h.highlightX+h.x0-h.x2*2), float32(h.highlightY+h.y0*3))
	h.path.LineTo(float32(h.highlightX+h.x0-h.x2*2), float32(h.highlightY+h.y0*4))
	h.path.LineTo(float32(h.highlightX+h.x0+h.x2*2), float32(h.highlightY+h.y0*4))
	h.path.LineTo(float32(h.highlightX+h.x0+h.x2*2), float32(h.highlightY+h.y0*3))
	h.path.LineTo(float32(h.highlightX+h.x0+h.x1), float32(h.highlightY+h.y0*3))
	h.path.Close()
	r, g, b, a := h.background[SquareColor(h.cursorRank, h.cursorFile)].RGBA()
	h.vs, h.is = h.path.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range h.vs {
		h.vs[i].SrcX = 1
		h.vs[i].SrcY = 1
		h.vs[i].ColorR = float32(r) / 0xffff
		h.vs[i].ColorG = float32(g) / 0xffff
		h.vs[i].ColorB = float32(b) / 0xffff
		h.vs[i].ColorA = float32(a) / 0xffff
	}
}

func (h *EnPassant) Draw(dst *ebiten.Image) {
	if h.visible {
		//dst.DrawTriangles(h.vs, h.is, whiteSubImage, h.op)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(h.highlightX), float64(h.highlightY))
		op.ColorScale.Scale(.2, .2, .2, 0.1)
		graphics.GetPiece(PiecePawn+PlayerBlack).Draw(dst, op)
	}
}

func (h *EnPassant) Debug(screen *ebiten.Image, debugX [8]int, debugY int) {
	if h.visible {

	}
}
