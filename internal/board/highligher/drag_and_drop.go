package highligher

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	"us.figge.chess/internal/board/graphics"
)

type DragHighlighter interface {
	Highlighter
	DragBegin(index uint8, pieceType uint8)
	DragOver(index uint8, pieceType uint8)
	DragEnd(from, to uint8, pieceType uint8, cancelled bool)
}

type DragAndDrop struct {
	*Highlight
	highlighter DragHighlighter
	background  color.Color
	dragColor   color.Color
	draggingOp  *ebiten.DrawImageOptions
	dragOffsetX int
	dragOffsetY int
	dragIndex   uint8
	dragOver    uint8
	dragPiece   *graphics.Piece
	dragging    bool
}

func NewDragAndDrop(highlighter DragHighlighter, squareSize int, dragColor, background color.Color) *DragAndDrop {
	dd := &DragAndDrop{
		Highlight:   NewHighlight(highlighter, squareSize, background),
		highlighter: highlighter,
		background:  background,
		dragColor:   dragColor,
		draggingOp:  &ebiten.DrawImageOptions{},
	}
	return dd
}

func (d *DragAndDrop) Update(x, y int) bool {
	changed := d.Highlight.Update(x, y)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && d.visible && d.piece != nil {
		d.dragPiece = d.piece
		d.dragOffsetX = d.highlightX - x
		d.dragOffsetY = d.highlightY - y
		d.draggingOp.GeoM.Reset()
		d.draggingOp.GeoM.Translate(float64(x+d.dragOffsetX), float64(y+d.dragOffsetY))
		d.dragIndex = d.index
		d.dragOver = d.index
		d.Highlight.background = d.dragColor
		d.dragging = true
		changed = true
		d.highlighter.DragBegin(d.dragOver, d.dragPiece.Type())
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		d.dragging = false
		changed = true
		d.Highlight.background = d.background
		d.highlighter.DragEnd(d.dragIndex, d.dragOver, d.dragPiece.Type(), d.dragOver == d.dragIndex || !d.visible)
	} else if d.dragging {
		d.draggingOp.GeoM.Reset()
		d.draggingOp.GeoM.Translate(float64(x+d.dragOffsetX), float64(y+d.dragOffsetY))
		moved := d.dragOver != d.Highlight.index
		d.dragOver = d.Highlight.index
		if moved {
			d.highlighter.DragOver(d.dragOver, d.dragPiece.Type())
		}
	}
	return changed
}

func (d *DragAndDrop) Hide() {
	d.CancelDrag()
	d.Highlight.Hide()
}

func (d *DragAndDrop) IsDragging() bool {
	return d.dragging
}

func (d *DragAndDrop) DragIndex() uint8 {
	return d.dragIndex
}

func (d *DragAndDrop) DrawDrag(screen *ebiten.Image) {
	if d.dragging && d.visible {
		d.dragPiece.Draw(screen, d.draggingOp)
	}
}

func (d *DragAndDrop) CancelDrag() {
	d.dragging = false
	d.Highlight.background = d.background
	d.highlighter.DragEnd(d.cursorRank, d.cursorFile, d.dragPiece.Type(), true)
}
