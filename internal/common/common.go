package common

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

type Configuration interface {
	SquareSize() int
	TranslateRFtoXY(rank, file int) (float64, float64)
	TranslateXYtoRF(x, y int) (int, int, bool)
	TranslateRFtoIndex(rank, file int) int
	TranslateIndexToRF(index int) (int, int)
	TranslateIndexToXY(index int) (float64, float64)
	TranslateRFtoN(rank, file int) string
	TranslateNtoRF(notation string) (int, int, bool)
	TranslateNtoIndex(notation string) (int, bool)

	EnableDebug() bool
	DebugX(file int) int
	DebugY() int

	Token(pieceType uint8) Token
	SheetImageSize() int
	ColorWhite() color.Color
	ColorBlack() color.Color
	ColorValid() color.Color
	ColorInvalid() color.Color
	ColorHighlight() color.Color
	ColorStrength() color.Color

	Turn(uint8) string
	HighlightAttacks() bool
	ShowStrength() bool
	ShowLabels() bool
	FontHeight() int

	TextAt(dst *ebiten.Image, str string, x, y int, size float64, color color.Color)
	TextSize(str string, size float64) (float64, float64)
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
	White uint8 = 0b00000000
	Black uint8 = 0b00000001

	Pawn   uint8 = 0b0000000
	Knight uint8 = 0b0000010
	Bishop uint8 = 0b0000100
	Rook   uint8 = 0b0000110
	Queen  uint8 = 0b0001000
	King   uint8 = 0b0001010

	CastleWQ uint8 = 0b00000010
	CastleWK uint8 = 0b00000100
	CastleBQ uint8 = 0b00001000
	CastleBK uint8 = 0b00010000
	EnPas
)

var FenPieceMap = map[byte]uint8{
	'p': Pawn | Black,
	'n': Knight | Black,
	'b': Bishop | Black,
	'r': Rook | Black,
	'q': Queen | Black,
	'k': King | Black,
	'P': Pawn | White,
	'N': Knight | White,
	'B': Bishop | White,
	'R': Rook | White,
	'Q': Queen | White,
	'K': King | White,
}
