package colors

import "image/color"

type Colors struct {
	black       color.Color
	playerWhite color.Color
	playerBlack color.Color
	valid       color.Color
	invalid     color.Color
	highlight   color.Color
}

func NewColors() *Colors {
	return &Colors{
		black:       &color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff},
		playerWhite: &color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		playerBlack: &color.RGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xff},
		valid:       &color.RGBA{R: 0x44, G: 0xff, B: 0x44, A: 0xff},
		invalid:     &color.RGBA{R: 0xff, G: 0x44, B: 0x44, A: 0xff},
		highlight:   &color.RGBA{R: 0x88, G: 0x22, B: 0x22, A: 0xff},
	}
}

func (c *Colors) Black() color.Color {
	return c.black
}
func (c *Colors) PlayerWhite() color.Color {
	return c.playerWhite
}
func (c *Colors) PlayerBlack() color.Color {
	return c.playerBlack
}
func (c *Colors) Valid() color.Color {
	return c.valid
}
func (c *Colors) Invalid() color.Color {
	return c.invalid
}
func (c *Colors) Highlight() color.Color {
	return c.highlight
}
func (c *Colors) Background() color.Color {
	return &color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}
}
func (c *Colors) Foreground() color.Color {
	return &color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
}
func (c *Colors) SetPlayerWhite(newColor *color.RGBA) {
	c.playerWhite = newColor
}
func (c *Colors) SetPlayerBlack(newColor *color.RGBA) {
	c.playerBlack = newColor
}
func (c *Colors) SetValid(newColor *color.RGBA) {
	c.valid = newColor
}
func (c *Colors) SetInvalid(newColor *color.RGBA) {
	c.invalid = newColor
}
func (c *Colors) SetHighlight(newColor *color.RGBA) {
	c.highlight = newColor
}
