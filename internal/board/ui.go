package board

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
	"strconv"
)

type ui struct {
	composite        *ebiten.Image
	background       *ebiten.Image
	labelingX        *ebiten.Image
	labelingY        *ebiten.Image
	labelingXOp      *ebiten.DrawImageOptions
	labelingFontSize float64
	foreground       *ebiten.Image
	strength         *ebiten.Image
	strengthOp       *ebiten.DrawImageOptions
	//highlightSquare  *square
}

func (b *Board) initializeImages() {
	b.labelingFontSize = float64(b.SquareSize()) * .17
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
	b.drawPieces(b.foreground)
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
}
