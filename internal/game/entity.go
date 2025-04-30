package game

import (
	"embed"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	_ "image/png"
	"us.figge.chess/internal/shared"
)

type Entity struct {
	img       *ebiten.Image
	tokenType uint8
	pieceType uint8
	colorType uint8
}

const (
	TokenWhitePawn   uint8 = 0b10101001
	TokenWhiteKnight uint8 = 0b01101010
	TokenWhiteBishop uint8 = 0b01001011
	TokenWhiteRook   uint8 = 0b10001100
	TokenWhiteQueen  uint8 = 0b00101101
	TokenWhiteKing   uint8 = 0b00001110
	TokenBlackPawn   uint8 = 0b10110001
	TokenBlackKnight uint8 = 0b01110010
	TokenBlackBishop uint8 = 0b01010011
	TokenBlackRook   uint8 = 0b10010100
	TokenBlackQueen  uint8 = 0b00110101
	TokenBlackKing   uint8 = 0b00010110
)

var (
	//go:embed assets/*.png
	assets embed.FS
	sheet  = mustLoadImage("assets/pieces.png")
)

func (p *Entity) Draw(dst *ebiten.Image, op *ebiten.DrawImageOptions) {
	dst.DrawImage(p.img, op)
}

func makeEntities(c shared.Configuration) map[uint8]*Entity {
	return map[uint8]*Entity{
		TokenWhitePawn & 0b00011111:   makeEntity(c, TokenWhitePawn),
		TokenWhiteKnight & 0b00011111: makeEntity(c, TokenWhiteKnight),
		TokenWhiteBishop & 0b00011111: makeEntity(c, TokenWhiteBishop),
		TokenWhiteRook & 0b00011111:   makeEntity(c, TokenWhiteRook),
		TokenWhiteQueen & 0b00011111:  makeEntity(c, TokenWhiteQueen),
		TokenWhiteKing & 0b00011111:   makeEntity(c, TokenWhiteKing),
		TokenBlackPawn & 0b00011111:   makeEntity(c, TokenBlackPawn),
		TokenBlackKnight & 0b00011111: makeEntity(c, TokenBlackKnight),
		TokenBlackBishop & 0b00011111: makeEntity(c, TokenBlackBishop),
		TokenBlackRook & 0b00011111:   makeEntity(c, TokenBlackRook),
		TokenBlackQueen & 0b00011111:  makeEntity(c, TokenBlackQueen),
		TokenBlackKing & 0b00011111:   makeEntity(c, TokenBlackKing),
	}
}
func makeEntity(c shared.Configuration, tokenType uint8) *Entity {
	entity := &Entity{
		img:       ebiten.NewImage(int(c.SquareSize()), int(c.SquareSize())),
		tokenType: tokenType,
		pieceType: tokenType & 0b00000111,
		colorType: tokenType & 0b00011000,
	}
	x := int(tokenType >> 5)
	y := int((tokenType >> 4) & 0b00000001)
	scale := float64(c.SquareSize()) / float64(c.SheetImageSize())
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	size := c.SheetImageSize()
	entity.img.DrawImage(sheet.SubImage(image.Rect(x*size, y*size, (x+1)*size, (y+1)*size)).(*ebiten.Image), op)

	return entity
}
func mustLoadImage(name string) *ebiten.Image {
	f, err := assets.Open(name)
	if err != nil {
		panic(err)
	}
	defer func() { _ = f.Close() }()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	return ebiten.NewImageFromImage(img)
}
