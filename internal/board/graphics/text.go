package graphics

import (
	"bytes"
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/text/language"
	"image/color"
)

var (
	//go:embed assets/Montserrat-Medium.ttf
	montserratTTF        []byte
	montserratFaceSource *text.GoTextFaceSource
)

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(montserratTTF))
	if err != nil {
		panic(err)
	}
	montserratFaceSource = s
}

func TextAt(dst *ebiten.Image, str string, x, y int, size float64, clr color.Color) {
	face := &text.GoTextFace{
		Source:    montserratFaceSource,
		Direction: text.DirectionLeftToRight,
		Size:      size,
		Language:  language.AmericanEnglish,
	}
	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(clr)
	op.GeoM.Translate(float64(x), float64(y))
	text.Draw(dst, str, face, op)
}

func TextSize(str string, size float64) (float64, float64) {
	face := &text.GoTextFace{
		Source:    montserratFaceSource,
		Direction: text.DirectionLeftToRight,
		Size:      size,
		Language:  language.AmericanEnglish,
	}
	return text.Measure(str, face, 0)
}
