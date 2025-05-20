package colors

import (
	"image/color"
)

type Colors struct {
	black       color.Color
	playerWhite color.Color
	playerBlack color.Color
	valid       color.Color
	invalid     color.Color
	highlight   color.Color
	dragStart   color.Color
	enPassant   color.Color
}

func NewColors() *Colors {
	return &Colors{
		black:       &color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff},
		playerWhite: &color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		playerBlack: &color.RGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xff},
		valid:       &color.RGBA{R: 0x00, G: 0x55, B: 0x00, A: 0xcc},
		//valid:     &color.RGBA{R: 0x55, G: 0x55, B: 0x55, A: 0xcc},
		invalid:   &color.RGBA{R: 0xff, G: 0x44, B: 0x44, A: 0xff},
		highlight: &color.RGBA{R: 0x70, G: 0x18, B: 0x18, A: 0x0},
		dragStart: &color.RGBA{R: 0xff, G: 0xff, B: 0x00, A: 0x80},
		enPassant: &color.RGBA{R: 0x00, G: 0xff, B: 0xff, A: 0x80},
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
func (c *Colors) DragStart() color.Color {
	return c.dragStart
}
func (c *Colors) EnPassant() color.Color {
	return c.enPassant
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
func (c *Colors) SetDragStart(newColor *color.RGBA) {
	c.dragStart = newColor
}
func (c *Colors) SetEnPassant(newColor *color.RGBA) {
	c.enPassant = newColor
}

func (c *Colors) Tints(tint color.Color) [2]color.Color {
	return [2]color.Color{
		tintColor(c.playerWhite, tint),
		tintColor(c.playerBlack, tint),
	}
}

func tintColor(base, tint color.Color) color.Color {
	r, g, b, a := base.RGBA()
	tr, tg, tb, ta := tint.RGBA()
	t := float32(ta&0xff) / 0xff
	ti := 1 - t
	return &color.RGBA{
		R: uint8(float32(r&0xff)*t + float32(tr&0xff)*ti),
		G: uint8(float32(g&0xff)*t + float32(tg&0xff)*ti),
		B: uint8(float32(b&0xff)*t + float32(tb&0xff)*ti),
		A: uint8(a),
	}
}
