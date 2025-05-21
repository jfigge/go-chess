package highlighers

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	. "us.figge.chess/internal/common"
)

type ValidMove struct {
	*Highlight
}

func NewValidMove(highlighter Highlighter, squareSize int, validColor [2]color.Color, index uint8) *ValidMove {
	vm := &ValidMove{
		Highlight: NewHighlight(highlighter, squareSize, validColor),
	}
	vm.UpdateByIndex(index)
	return vm
}

func (h *ValidMove) UpdateByIndex(index uint8) {
	h.Highlight.UpdateByIndex(index)
	h.highlightX += h.squareSize / 2
	h.highlightY += h.squareSize / 2
}

func (h *ValidMove) Draw(dst *ebiten.Image) {
	if h.visible {
		vector.DrawFilledCircle(
			dst,
			float32(h.highlightX),
			float32(h.highlightY),
			float32(h.squareSize/5),
			h.background[SquareColor(h.cursorRank, h.cursorFile)], false,
		)
	}
}

func (h *ValidMove) Debug(screen *ebiten.Image, debugX [8]int, debugY int) {
	if h.visible {

	}
}
