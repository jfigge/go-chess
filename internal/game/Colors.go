package game

import "image/color"

func OptWhiteRGB(red, green, blue uint8) GameOptions {
	return func(g *Game) {
		g.ColorScheme.white = &color.RGBA{R: red, G: green, B: blue, A: 0xff}
	}
}
func OptBlackRGB(red, green, blue uint8) GameOptions {
	return func(g *Game) {
		g.ColorScheme.black = &color.RGBA{R: red, G: green, B: blue, A: 0xff}
	}
}
func OptValidRGBA(red, green, blue, alpha uint8) GameOptions {
	return func(g *Game) {
		g.ColorScheme.valid = &color.RGBA{R: red, G: green, B: blue, A: alpha}
	}
}
func OptInvalidRGBA(red, green, blue, alpha uint8) GameOptions {
	return func(g *Game) {
		g.ColorScheme.invalid = &color.RGBA{R: red, G: green, B: blue, A: alpha}
	}
}
func OptHighlightRGBA(red, green, blue, alpha uint8) GameOptions {
	return func(g *Game) {
		g.ColorScheme.highlight = &color.RGBA{R: red, G: green, B: blue, A: alpha}
	}
}

type ColorScheme struct {
	white     color.Color
	black     color.Color
	valid     color.Color
	invalid   color.Color
	highlight color.Color
}

func newColorScheme() *ColorScheme {
	return &ColorScheme{
		white:     &color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		black:     &color.RGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xff},
		valid:     &color.RGBA{R: 0x44, G: 0xff, B: 0x44, A: 0xff},
		invalid:   &color.RGBA{R: 0xff, G: 0x44, B: 0x44, A: 0xff},
		highlight: &color.RGBA{R: 0x88, G: 0x22, B: 0x22, A: 0xff},
	}
}

func (c *ColorScheme) ColorWhite() color.Color {
	return c.white
}
func (c *ColorScheme) ColorBlack() color.Color {
	return c.black
}
func (c *ColorScheme) ColorValid() color.Color {
	return c.valid
}
func (c *ColorScheme) ColorInvalid() color.Color {
	return c.invalid
}
func (c *ColorScheme) ColorHighlight() color.Color {
	return c.highlight
}
