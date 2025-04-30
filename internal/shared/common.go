package shared

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

type Configuration interface {
	SquareSize() int
	WhiteColor() color.Color
	BlackColor() color.Color
	Translate(rank, file int) (float64, float64)
	Token(pieceType uint) Token
	SheetImageSize() int
}

type Token interface {
	Draw(dst *ebiten.Image, op *ebiten.DrawImageOptions)
}
