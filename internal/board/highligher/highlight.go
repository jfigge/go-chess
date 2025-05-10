package highligher

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"strconv"
	"us.figge.chess/internal/board/graphics"
	. "us.figge.chess/internal/common"
)

type Highlighter interface {
	GetPieceType(rank, file uint8) (uint8, bool)
}

type Highlight struct {
	highlighter Highlighter
	squareSize  int
	background  color.Color
	visible     bool
	index       uint8
	cursorX     int
	cursorY     int
	cursorRank  uint8
	cursorFile  uint8
	highlightX  float32
	highlightY  float32
	piece       *graphics.Piece
}

func NewHighlight(highlighter Highlighter, squareSize int, background color.Color) *Highlight {
	h := &Highlight{
		index:       0xff,
		highlighter: highlighter,
		squareSize:  squareSize,
		background:  background,
	}
	h.Update(-1, -1)
	return h
}
func (h *Highlight) Update(x, y int) bool {
	rank, file, inRange := XYtoRF(x, y, h.squareSize)
	changed := h.visible != inRange
	if !inRange {
		h.visible = false
		return changed
	}
	h.visible = true
	index := RFtoI(rank, file)
	if index == h.index {
		return changed
	}
	hx, hy := RFtoXY(rank, file, h.squareSize)
	h.index = index
	h.cursorX = x
	h.cursorY = y
	h.cursorRank = rank
	h.cursorFile = file
	h.highlightX = hx
	h.highlightY = hy

	pieceType, present := h.highlighter.GetPieceType(rank, file)
	if present {
		h.piece = graphics.GetPiece(pieceType)
	} else {
		h.piece = nil
	}
	return true
}

func (h *Highlight) Draw(dst *ebiten.Image) {
	if h.visible {
		vector.DrawFilledRect(dst, h.highlightX, h.highlightY, float32(h.squareSize), float32(h.squareSize), h.background, false)
	}
}

func (h *Highlight) Debug(screen *ebiten.Image, debugX [8]int, debugY int) {
	if h.visible {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("R,F: %d,%d", h.cursorRank, h.cursorFile), debugX[1], debugY)
		ebitenutil.DebugPrintAt(screen, "Index: "+strconv.Itoa(int(h.index)), debugX[0], debugY)
		if h.piece != nil {
			ebitenutil.DebugPrintAt(screen, h.piece.ColorName()+" "+h.piece.Name(), debugX[4], debugY)
		}
	}
}
