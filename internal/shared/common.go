package shared

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

type Configuration interface {
	SquareSize() uint
	WhiteColor() color.Color
	BlackColor() color.Color
	CastleColor() color.Color
	Translate(rank, file uint8) (float64, float64)
	Token(pieceType uint8) Token
	SheetImageSize() int
}

type Token interface {
	Draw(target *ebiten.Image, op *ebiten.DrawImageOptions)
}

const (
	White uint8 = 0b00001000
	Black uint8 = 0b00010000

	Pawn   uint8 = 0b00000001
	Knight uint8 = 0b00000010
	Bishop uint8 = 0b00000011
	Rook   uint8 = 0b00000100
	Queen  uint8 = 0b00000101
	King   uint8 = 0b00000110
)

var FenPieceMap = map[byte]uint8{
	'p': Pawn | Black,
	'r': Rook | Black,
	'n': Knight | Black,
	'b': Bishop | Black,
	'q': Queen | Black,
	'k': King | Black,
	'P': Pawn | White,
	'R': Rook | White,
	'N': Knight | White,
	'B': Bishop | White,
	'Q': Queen | White,
	'K': King | White,
}
