package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRFtoI(t *testing.T) {
	tests := map[string]struct {
		index uint8
		rank  uint8
		file  uint8
	}{
		"top-left":     {0, 8, 1},
		"top-right":    {7, 8, 8},
		"bottom-left":  {56, 1, 1},
		"bottom-right": {63, 1, 8},
	}
	for name, test := range tests {
		t.Run(name, func(tt *testing.T) {
			index := RFtoI(test.rank, test.file)
			assert.Equal(tt, test.index, index)
		})
	}
}

func TestItoRF(t *testing.T) {
	tests := map[string]struct {
		index uint8
		rank  uint8
		file  uint8
	}{
		"top-left":     {0, 8, 1},
		"top-right":    {7, 8, 8},
		"bottom-left":  {56, 1, 1},
		"bottom-right": {63, 1, 8},
	}
	for name, test := range tests {
		t.Run(name, func(tt *testing.T) {
			rank, file := ItoRF(test.index)
			assert.Equal(tt, test.rank, rank)
			assert.Equal(tt, test.file, file)
		})
	}
}

func TestRFtoB(t *testing.T) {
	tests := map[string]struct {
		bit  uint64
		rank uint8
		file uint8
	}{
		"top-left":     {1 << 63, 8, 1},
		"top-right":    {1 << 56, 8, 8},
		"bottom-left":  {1 << 7, 1, 1},
		"bottom-right": {1 << 0, 1, 8},
	}
	for name, test := range tests {
		t.Run(name, func(tt *testing.T) {
			bit := RFtoB(test.rank, test.file)
			assert.Equal(tt, test.bit, bit)
		})
	}
}

func TestBtoRF(t *testing.T) {
	tests := map[string]struct {
		bit  uint64
		rank uint8
		file uint8
	}{
		"top-left":     {1 << 63, 8, 1},
		"top-right":    {1 << 56, 8, 8},
		"bottom-left":  {1 << 7, 1, 1},
		"bottom-right": {1 << 0, 1, 8},
	}
	for name, test := range tests {
		t.Run(name, func(tt *testing.T) {
			rank, file := BtoRF(test.bit)
			assert.Equal(tt, test.rank, rank)
			assert.Equal(tt, test.file, file)
		})
	}
}

func TestNtoRF(t *testing.T) {
	tests := map[string]struct {
		n    string
		rank uint8
		file uint8
		ok   bool
	}{
		"a1": {"a1", 1, 1, true},
		"a8": {"a8", 8, 1, true},
		"h1": {"h1", 1, 8, true},
		"h8": {"h8", 8, 8, true},
	}
	for name, test := range tests {
		t.Run(name, func(tt *testing.T) {
			rank, file, valid := NtoRF(test.n)
			assert.Equal(tt, test.rank, rank)
			assert.Equal(tt, test.file, file)
			assert.Equal(tt, test.ok, valid)
		})
	}
}

func TestRFtoN(t *testing.T) {
	tests := map[string]struct {
		n    string
		rank uint8
		file uint8
		ok   bool
	}{
		"a1": {"a1", 1, 1, true},
		"a8": {"a8", 8, 1, true},
		"h1": {"h1", 1, 8, true},
		"h8": {"h8", 8, 8, true},
	}
	for name, test := range tests {
		t.Run(name, func(tt *testing.T) {
			n := RFtoN(test.rank, test.file)
			assert.Equal(tt, test.n, n)
		})
	}
}

func TestBtoI(t *testing.T) {
	tests := map[string]struct {
		index uint8
		bit   uint64
	}{
		"top-left":     {0, 1 << 63},
		"top-right":    {7, 1 << 56},
		"bottom-left":  {56, 1 << 7},
		"bottom-right": {63, 1 << 0},
	}
	for name, test := range tests {
		t.Run(name, func(tt *testing.T) {
			bit := BtoI(test.bit)
			assert.Equal(tt, test.index, bit)
		})
	}
}

func TestItoB(t *testing.T) {
	tests := map[string]struct {
		index uint8
		bit   uint64
	}{
		"top-left":     {0, 1 << 63},
		"top-right":    {7, 1 << 56},
		"bottom-left":  {56, 1 << 7},
		"bottom-right": {63, 1 << 0},
	}
	for name, test := range tests {
		t.Run(name, func(tt *testing.T) {
			bit := ItoB(test.index)
			assert.Equal(tt, test.bit, bit)
		})
	}
}

func TestPTtoBB(t *testing.T) {
	tests := map[string]struct {
		pieceType   uint8
		pieceBoard  uint8
		playerBoard uint8
	}{
		"white pawn":   {PiecePawn | PlayerWhite, BitPawns, BitWhite},
		"white knight": {PieceKnight | PlayerWhite, BitKnights, BitWhite},
		"white bishop": {PieceBishop | PlayerWhite, BitBishops, BitWhite},
		"white rook":   {PieceRook | PlayerWhite, BitRooks, BitWhite},
		"white queen":  {PieceQueen | PlayerWhite, BitQueens, BitWhite},
		"white king":   {PieceKing | PlayerWhite, BitKings, BitWhite},
		"black pawn":   {PiecePawn | PlayerBlack, BitPawns, BitBlack},
		"black knight": {PieceKnight | PlayerBlack, BitKnights, BitBlack},
		"black bishop": {PieceBishop | PlayerBlack, BitBishops, BitBlack},
		"black rook":   {PieceRook | PlayerBlack, BitRooks, BitBlack},
		"black queen":  {PieceQueen | PlayerBlack, BitQueens, BitBlack},
		"black king":   {PieceKing | PlayerBlack, BitKings, BitBlack},
	}
	for name, test := range tests {
		t.Run(name, func(tt *testing.T) {
			pieceBoard, playerBoard := PTtoBB(test.pieceType)
			assert.Equal(tt, test.pieceBoard, pieceBoard)
			assert.Equal(tt, test.playerBoard, playerBoard)
		})
	}
}
