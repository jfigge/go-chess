package game

import (
	"embed"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	_ "image/png"
	"strings"
	"us.figge.chess/internal/shared"
)

type Entity struct {
	img       *ebiten.Image
	name      string
	color     string
	fen       string
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

func (e *Entity) Draw(dst *ebiten.Image, op *ebiten.DrawImageOptions) {
	dst.DrawImage(e.img, op)
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
	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
	op.GeoM.Scale(scale, scale)
	size := c.SheetImageSize()
	entity.img.DrawImage(sheet.SubImage(image.Rect(x*size, y*size, (x+1)*size, (y+1)*size)).(*ebiten.Image), op)
	switch entity.pieceType {
	case shared.Pawn:
		entity.name = "Pawn"
		entity.fen = "p"
	case shared.Knight:
		entity.name = "Knight"
		entity.fen = "n"
	case shared.Bishop:
		entity.name = "Bishop"
		entity.fen = "b"
	case shared.Rook:
		entity.name = "Rook"
		entity.fen = "r"
	case shared.Queen:
		entity.name = "Queen"
		entity.fen = "q"
	case shared.King:
		entity.name = "King"
		entity.fen = "k"
	}
	switch entity.colorType {
	case shared.White:
		entity.color = "White"
		entity.fen = strings.ToUpper(entity.fen)
	case shared.Black:
		entity.color = "Black"
		entity.fen = strings.ToLower(entity.fen)
	}
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

func (e *Entity) IsPawn() bool   { return e.pieceType == shared.Pawn }
func (e *Entity) IsKnight() bool { return e.pieceType == shared.Knight }
func (e *Entity) IsBishop() bool { return e.pieceType == shared.Bishop }
func (e *Entity) IsRook() bool   { return e.pieceType == shared.Rook }
func (e *Entity) IsQueen() bool  { return e.pieceType == shared.Queen }
func (e *Entity) IsKing() bool   { return e.pieceType == shared.King }
func (e *Entity) IsWhite() bool  { return e.colorType == shared.White }
func (e *Entity) IsBlack() bool  { return e.colorType == shared.Black }
func (e *Entity) Name() string {
	return e.name
}
func (e *Entity) Color() string {
	return e.color
}
func (e *Entity) Fen() byte {
	return e.fen[0]
}
