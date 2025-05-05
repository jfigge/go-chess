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

	EnableDebug() bool
	DebugX(file int) int
	DebugY() int

	Piece(pieceType uint8) Piece
	SheetImageSize() int
	ColorWhite() color.Color
	ColorBlack() color.Color
	ColorValid() color.Color
	ColorInvalid() color.Color
	ColorHighlight() color.Color
	ColorStrength() color.Color

	TurnName(uint8) string
	HighlightAttacks() bool
	ShowStrength() bool
	ShowLabels() bool
	FontHeight() int

	TextAt(dst *ebiten.Image, str string, x, y int, size float64, color color.Color)
	TextSize(str string, size float64) (float64, float64)
}

type Piece interface {
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
	// Move stored in White Pawn Rank 1
	TurnMask uint8 = 0b00000001
	White    uint8 = 0b00000000
	Black    uint8 = 0b00000001

	PieceMask uint8 = 0b00001110
	Pawn      uint8 = 0b00000000
	Knight    uint8 = 0b00000010
	Bishop    uint8 = 0b00000100
	Rook      uint8 = 0b00000110
	Queen     uint8 = 0b00001000
	King      uint8 = 0b00001010

	// Castling rights stored in White Pawn Rank 1
	CastleRightsMask       uint8 = 0b00011110
	CastleRightsWhiteKing  uint8 = 0b00000010
	CastleRightsWhiteQueen uint8 = 0b00000100
	CastleRightsBlackKing  uint8 = 0b00001000
	CastleRightsBlackQueen uint8 = 0b00010000

	// EnPassant stored in Black Pawn Rank 8
	EnPassantMask uint8 = 0b11111111
	EnPassantA    uint8 = 0b00001000
	EnPassantB    uint8 = 0b00000111
	EnPassantC    uint8 = 0b00000110
	EnPassantD    uint8 = 0b00000101
	EnPassantE    uint8 = 0b00000100
	EnPassantF    uint8 = 0b00000011
	EnPassantG    uint8 = 0b00000010
	EnPassantH    uint8 = 0b00000001

	// Bit board index
	WhitePawn   uint8 = 0b00000000
	BlackPawn   uint8 = 0b00000001
	WhiteKnight uint8 = 0b00000010
	BlackKnight uint8 = 0b00000011
	WhiteBishop uint8 = 0b00000100
	BlackBishop uint8 = 0b00000101
	WhiteRook   uint8 = 0b00000110
	BlackRook   uint8 = 0b00000111
	WhiteQueen  uint8 = 0b00001000
	BlackQueen  uint8 = 0b00001001
	WhiteKing   uint8 = 0b00001010
	BlackKing   uint8 = 0b00001011
)
