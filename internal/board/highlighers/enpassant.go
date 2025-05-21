package highlighers

import (
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
}

func NewEnPassant(highlighter Highlighter, squareSize int, validColor [2]color.Color) *EnPassant {
	ep := &EnPassant{
		Highlight: NewHighlight(highlighter, squareSize, validColor),
	}
	return ep
}

func (h *EnPassant) UpdateByIndex(index uint8) {
	h.Highlight.UpdateByIndex(index)
	player := PlayerWhite
	if index < 32 {
		player = PlayerBlack
	}
	h.piece = graphics.GetPiece(PiecePawn + player)
}

func (h *EnPassant) Draw(dst *ebiten.Image) {
	if h.visible {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(h.highlightX), float64(h.highlightY))
		op.ColorScale.Scale(.25, .25, .25, 0.2)
		h.piece.Draw(dst, op)
	}
}

func (h *EnPassant) Debug(screen *ebiten.Image, debugX [8]int, debugY int) {
	if h.visible {

	}
}
