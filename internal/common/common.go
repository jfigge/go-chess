package common

import (
	"math"
)

const (
	boardPattern uint64 = 6172840429334713770
	notation            = "abcdefgh"
)

const (
	PlayerMask  uint8 = 0b00000001
	PlayerWhite uint8 = 0b00000000
	PlayerBlack uint8 = 0b00000001
)

// Piece Types - 0 - 5 shifted left by 1 to allow for the color
const (
	PieceMask   uint8 = 0b00001110
	PiecePawn   uint8 = 0b00000000
	PieceKnight uint8 = 0b00000010
	PieceBishop uint8 = 0b00000100
	PieceRook   uint8 = 0b00000110
	PieceQueen  uint8 = 0b00001000
	PieceKing   uint8 = 0b00001010
)

// Castling rights stored in status
const (
	CastleRightsMask       uint8 = 0b00011110
	CastleRightsWhiteMask  uint8 = 0b00000110
	CastleRightsWhiteKing  uint8 = 0b00000010
	CastleRightsWhiteQueen uint8 = 0b00000100
	CastleRightsBlackMask  uint8 = 0b00011000
	CastleRightsBlackKing  uint8 = 0b00001000
	CastleRightsBlackQueen uint8 = 0b00010000
)

// Bit Boards
const (
	BitWhite     uint8 = 0b00000000 // 0
	BitBlack     uint8 = 0b00000001 // 1
	BitPawns     uint8 = 0b00000010 // 2  0
	BitKnights   uint8 = 0b00000011 // 3  1
	BitBishops   uint8 = 0b00000100 // 4  2
	BitRooks     uint8 = 0b00000101 // 5  3
	BitQueens    uint8 = 0b00000110 // 6  4
	BitKings     uint8 = 0b00000111 // 7  5
	BitEnPassant uint8 = 0b00001000 // 8  6
)

// RFtoI Rank and File to Index
func RFtoI(rank, file uint8) uint8 {
	return (8-rank)*8 + file - 1
}

// ItoRF Index to Rank and File
func ItoRF(index uint8) (uint8, uint8) {
	return 8 - index/8, index%8 + 1
}

// RFtoB Rank and File to bitboard bit
func RFtoB(rank, file uint8) uint64 {
	return ItoB(RFtoI(rank, file))
}

// ItoB Index to bitboard bit
func ItoB(index uint8) uint64 {
	return 1 << (63 - index)
}

// BtoI Bitboard bit to Index
func BtoI(bit uint64) uint8 {
	return 63 - uint8(math.Log(float64(bit))/math.Log(2))
}

// BtoRF Bitboard bit index to Rank and File
func BtoRF(bit uint64) (uint8, uint8) {
	return ItoRF(BtoI(bit))
}

// PTtoBB PieceType to piece and player Bitboards
func PTtoBB(pieceType uint8) (uint8, uint8) {
	return pieceType&PieceMask>>1 + 2, pieceType & PlayerMask
}

// NtoRF Notation to Rank and File
func NtoRF(n string) (uint8, uint8, bool) {
	if len(n) != 2 {
		return 0, 0, false
	}
	rank := n[1] - '0'
	file := n[0] - 'a' + 1
	if rank < 1 || rank > 8 || file < 1 || file > 8 {
		return 0, 0, false
	}
	return rank, file, true
}
func RFtoN(rank, file uint8) string {
	return string(notation[file-1]) + string(rank+'0')
}
func FtoN(file uint8) string {
	return string(notation[file-1])
}

func XYtoRF(x, y, squareSize int) (uint8, uint8, bool) {
	rank := uint8(8 - (y)/squareSize)
	file := uint8((x)/squareSize + 1)
	if x < 1 || y < 3 || rank < 1 || rank > 8 || file < 1 || file > 8 {
		return 0, 0, false
	}
	return rank, file, true
}

func RFtoXY(rank, file uint8, squareSize int) (int, int) {
	return int(file-1) * squareSize, int(8-rank) * squareSize
}

func SquareColor(rank, file uint8) uint64 {
	i := RFtoI(rank, file)
	return (boardPattern >> (63 - i)) & 1
}
