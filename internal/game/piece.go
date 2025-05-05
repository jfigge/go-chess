package game

import (
	_ "image/png"

	"embed"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"us.figge.chess/internal/common"
)

type Piece struct {
	img       *ebiten.Image
	pieceType uint8
	colorType uint8
}

const (
	fen              = "PpNnBbRrQqKk"
	pieceWhitePawn   = 0b01010000 | common.White | common.Pawn
	pieceWhiteKnight = 0b00110000 | common.White | common.Knight
	pieceWhiteBishop = 0b00100000 | common.White | common.Bishop
	pieceWhiteRook   = 0b01000000 | common.White | common.Rook
	pieceWhiteQueen  = 0b00010000 | common.White | common.Queen
	pieceWhiteKing   = 0b00000000 | common.White | common.King
	pieceBlackPawn   = 0b01010000 | common.Black | common.Pawn
	pieceBlackKnight = 0b00110000 | common.Black | common.Knight
	pieceBlackBishop = 0b00100000 | common.Black | common.Bishop
	pieceBlackRook   = 0b01000000 | common.Black | common.Rook
	pieceBlackQueen  = 0b00010000 | common.Black | common.Queen
	pieceBlackKing   = 0b00000000 | common.Black | common.King
)

var (
	//go:embed assets/*.png
	assets embed.FS
	sheet  = mustLoadImage("assets/pieces.png")
	names  = []string{"Pawn", "Knight", "Bishop", "Rook", "Queen", "King"}
	colors = []string{"White", "Black"}
)

func initialize(c common.Configuration) map[uint8]common.Piece {
	scale := float64(c.SquareSize()) / float64(c.SheetImageSize())
	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
	op.GeoM.Scale(scale, scale)
	return map[uint8]common.Piece{
		pieceWhitePawn & 0b00001111:   makeEntity(c, op, pieceWhitePawn),
		pieceWhiteKnight & 0b00001111: makeEntity(c, op, pieceWhiteKnight),
		pieceWhiteBishop & 0b00001111: makeEntity(c, op, pieceWhiteBishop),
		pieceWhiteRook & 0b00001111:   makeEntity(c, op, pieceWhiteRook),
		pieceWhiteQueen & 0b00001111:  makeEntity(c, op, pieceWhiteQueen),
		pieceWhiteKing & 0b00001111:   makeEntity(c, op, pieceWhiteKing),
		pieceBlackPawn & 0b00001111:   makeEntity(c, op, pieceBlackPawn),
		pieceBlackKnight & 0b00001111: makeEntity(c, op, pieceBlackKnight),
		pieceBlackBishop & 0b00001111: makeEntity(c, op, pieceBlackBishop),
		pieceBlackRook & 0b00001111:   makeEntity(c, op, pieceBlackRook),
		pieceBlackQueen & 0b00001111:  makeEntity(c, op, pieceBlackQueen),
		pieceBlackKing & 0b00001111:   makeEntity(c, op, pieceBlackKing),
	}
}

func (e *Piece) Draw(dst *ebiten.Image, op *ebiten.DrawImageOptions) {
	dst.DrawImage(e.img, op)
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

func makeEntity(c common.Configuration, op *ebiten.DrawImageOptions, tokenType uint8) *Piece {
	entity := &Piece{
		img:       ebiten.NewImage(c.SquareSize(), c.SquareSize()),
		pieceType: tokenType & common.PieceMask,
		colorType: tokenType & common.TurnMask,
	}
	x := int(tokenType >> 4)
	y := int(tokenType & common.TurnMask)
	size := c.SheetImageSize()
	entity.img.DrawImage(sheet.SubImage(image.Rect(x*size, y*size, (x+1)*size, (y+1)*size)).(*ebiten.Image), op)
	return entity
}

func (e *Piece) IsPawn() bool   { return e.pieceType == common.Pawn }
func (e *Piece) IsKnight() bool { return e.pieceType == common.Knight }
func (e *Piece) IsBishop() bool { return e.pieceType == common.Bishop }
func (e *Piece) IsRook() bool   { return e.pieceType == common.Rook }
func (e *Piece) IsQueen() bool  { return e.pieceType == common.Queen }
func (e *Piece) IsKing() bool   { return e.pieceType == common.King }
func (e *Piece) IsWhite() bool  { return e.colorType == common.White }
func (e *Piece) IsBlack() bool  { return e.colorType == common.Black }
func (e *Piece) Name() string   { return names[e.pieceType>>1] }
func (e *Piece) Color() string  { return colors[e.colorType] }
func (e *Piece) Fen() byte      { return fen[e.pieceType|e.colorType] }
