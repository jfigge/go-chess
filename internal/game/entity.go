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
	pieceType uint8
	colorType uint8
}

const (
	fen              = "pnbrqk"
	TokenWhitePawn   = 0b01010000 | shared.Pawn | shared.White
	TokenWhiteKnight = 0b00110000 | shared.Knight | shared.White
	TokenWhiteBishop = 0b00100000 | shared.Bishop | shared.White
	TokenWhiteRook   = 0b01000000 | shared.Rook | shared.White
	TokenWhiteQueen  = 0b00010000 | shared.Queen | shared.White
	TokenWhiteKing   = 0b00000000 | shared.King | shared.White
	TokenBlackPawn   = 0b01010000 | shared.Pawn | shared.Black
	TokenBlackKnight = 0b00110000 | shared.Knight | shared.Black
	TokenBlackBishop = 0b00100000 | shared.Bishop | shared.Black
	TokenBlackRook   = 0b01000000 | shared.Rook | shared.Black
	TokenBlackQueen  = 0b00010000 | shared.Queen | shared.Black
	TokenBlackKing   = 0b00000000 | shared.King | shared.Black
)

var (
	//go:embed assets/*.png
	assets embed.FS
	sheet  = mustLoadImage("assets/pieces.png")
	names  = []string{"Pawn", "Knight", "Bishop", "Rook", "Queen", "King"}
	colors = []string{"White", "Black"}
)

func (e *Entity) Draw(dst *ebiten.Image, op *ebiten.DrawImageOptions) {
	dst.DrawImage(e.img, op)
}

func makeEntities(c shared.Configuration) map[uint8]*Entity {
	return map[uint8]*Entity{
		TokenWhitePawn & 0b00001111:   makeEntity(c, TokenWhitePawn),
		TokenWhiteKnight & 0b00001111: makeEntity(c, TokenWhiteKnight),
		TokenWhiteBishop & 0b00001111: makeEntity(c, TokenWhiteBishop),
		TokenWhiteRook & 0b00001111:   makeEntity(c, TokenWhiteRook),
		TokenWhiteQueen & 0b00001111:  makeEntity(c, TokenWhiteQueen),
		TokenWhiteKing & 0b00001111:   makeEntity(c, TokenWhiteKing),
		TokenBlackPawn & 0b00001111:   makeEntity(c, TokenBlackPawn),
		TokenBlackKnight & 0b00001111: makeEntity(c, TokenBlackKnight),
		TokenBlackBishop & 0b00001111: makeEntity(c, TokenBlackBishop),
		TokenBlackRook & 0b00001111:   makeEntity(c, TokenBlackRook),
		TokenBlackQueen & 0b00001111:  makeEntity(c, TokenBlackQueen),
		TokenBlackKing & 0b00001111:   makeEntity(c, TokenBlackKing),
	}
}
func makeEntity(c shared.Configuration, tokenType uint8) *Entity {
	entity := &Entity{
		img:       ebiten.NewImage(c.SquareSize(), c.SquareSize()),
		pieceType: tokenType & 0b00001110,
		colorType: tokenType & 0b00000001,
	}
	x := int(tokenType >> 4)
	y := int(tokenType & 0b00000001)
	scale := float64(c.SquareSize()) / float64(c.SheetImageSize())
	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
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

func (e *Entity) IsPawn() bool   { return e.pieceType == shared.Pawn }
func (e *Entity) IsKnight() bool { return e.pieceType == shared.Knight }
func (e *Entity) IsBishop() bool { return e.pieceType == shared.Bishop }
func (e *Entity) IsRook() bool   { return e.pieceType == shared.Rook }
func (e *Entity) IsQueen() bool  { return e.pieceType == shared.Queen }
func (e *Entity) IsKing() bool   { return e.pieceType == shared.King }
func (e *Entity) IsWhite() bool  { return e.colorType == shared.White }
func (e *Entity) IsBlack() bool  { return e.colorType == shared.Black }
func (e *Entity) Name() string   { return names[e.pieceType>>1] }
func (e *Entity) Color() string  { return colors[e.colorType] }
func (e *Entity) Fen() byte      { return fen[e.pieceType>>1] }
