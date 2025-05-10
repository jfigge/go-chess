package board

import "image/color"

type Option func(b *Board)

func OptDebugEnabled(enabled bool) Option {
	return func(b *Board) {
		b.debugEnabled = enabled
	}
}

func OptSquareSize(size int) Option {
	return func(b *Board) {
		b.squareSize = size
	}
}

func OptWhiteRGB(red, green, blue uint8) Option {
	return func(b *Board) {
		b.colors.SetPlayerWhite(&color.RGBA{R: red, G: green, B: blue, A: 0xff})
	}
}
func OptBlackRGB(red, green, blue uint8) Option {
	return func(b *Board) {
		b.colors.SetPlayerBlack(&color.RGBA{R: red, G: green, B: blue, A: 0xff})
	}
}
func OptValidRGBA(red, green, blue, alpha uint8) Option {
	return func(b *Board) {
		b.colors.SetValid(&color.RGBA{R: red, G: green, B: blue, A: alpha})
	}
}
func OptInvalidRGBA(red, green, blue, alpha uint8) Option {
	return func(b *Board) {
		b.colors.SetInvalid(&color.RGBA{R: red, G: green, B: blue, A: alpha})
	}
}
func OptHighlightRGBA(red, green, blue, alpha uint8) Option {
	return func(b *Board) {
		b.colors.SetHighlight(&color.RGBA{R: red, G: green, B: blue, A: alpha})
	}
}
