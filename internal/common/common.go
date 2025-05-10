package common

const (
	notation = "abcdefgh"
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
	CastleRightsWhiteKing  uint8 = 0b00000010
	CastleRightsWhiteQueen uint8 = 0b00000100
	CastleRightsBlackKing  uint8 = 0b00001000
	CastleRightsBlackQueen uint8 = 0b00010000
)

// RFtoI Rank and File to Index
func RFtoI(rank, file uint8) uint8 {
	return (8-rank)*8 + file - 1
}

// ItoRF Index to Rank and File
func ItoRF(index uint8) (uint8, uint8) {
	return 8 - index/8, index%8 + 1
}

// RFtoB Rank and File to bitboard bit index
func RFtoB(rank, file uint8) uint8 {
	return (rank-1)*8 + 8 - file
}

// BtoRF Bitboard bit index to Rank and File
func BtoRF(index uint8) (uint8, uint8) {
	return index/8 + 1, 8 - index%8
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
