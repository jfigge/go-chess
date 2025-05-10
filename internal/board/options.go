package board

import "image/color"

type Option func(b *Board)

func OptSquareSize(size int) Option {
	return func(b *Board) {
		b.squareSize = size
	}
}

func OptFontHeight(height int) Option {
	return func(d *Board) {
		d.fontHeight = height
	}
}

func OptWhiteRGB(red, green, blue uint8) Option {
	return func(b *Board) {
		b.colorWhite = &color.RGBA{R: red, G: green, B: blue, A: 0xff}
	}
}
func OptBlackRGB(red, green, blue uint8) Option {
	return func(b *Board) {
		b.colorBlack = &color.RGBA{R: red, G: green, B: blue, A: 0xff}
	}
}
func OptValidRGBA(red, green, blue, alpha uint8) Option {
	return func(b *Board) {
		b.colorValid = &color.RGBA{R: red, G: green, B: blue, A: alpha}
	}
}
func OptInvalidRGBA(red, green, blue, alpha uint8) Option {
	return func(b *Board) {
		b.colorInvalid = &color.RGBA{R: red, G: green, B: blue, A: alpha}
	}
}
func OptHighlightRGBA(red, green, blue, alpha uint8) Option {
	return func(b *Board) {
		b.colorHighlight = &color.RGBA{R: red, G: green, B: blue, A: alpha}
	}
}
