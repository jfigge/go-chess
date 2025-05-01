package shared

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

type Configuration interface {
	SquareSize() uint
	TranslateRFtoXY(rank, file uint8) (float64, float64)
	TranslateXYtoRF(x, y int) (uint8, uint8, bool)
	TranslateRFtoIndex(rank, file uint8) uint8
	TranslateIndexToRF(index uint8) (uint8, uint8)
	TranslateIndexToXY(index uint8) (float64, float64)
	TranslateRFtoN(rank, file uint8) string
	TranslateNtoRF(notation string) (uint8, uint8, bool)
	TranslateNtoIndex(notation string) (uint8, bool)

	EnableDebug() bool
	DebugX(rank uint8) int
	DebugY() int
	DebugFen() int

	Token(pieceType uint8) Token
	SheetImageSize() int
	ColorWhite() color.Color
	ColorBlack() color.Color
	ColorValid() color.Color
	ColorInvalid() color.Color
	ColorHighlight() color.Color
}

type Token interface {
	Draw(target *ebiten.Image, op *ebiten.DrawImageOptions)
	Name() string
	Color() string
	Fen() byte
	IsPawn() bool
	IsKnight() bool
	IsBishop() bool
	IsRook() bool
	IsQueen() bool
	IsKing() bool
	IsWhite() bool
	IsBlack() bool
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
