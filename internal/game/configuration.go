package game

import (
	"strings"
	"us.figge.chess/internal/shared"
)

const (
	notation = "abcdefgh"
)

func (g *Game) SquareSize() uint {
	return g.squareSize
}

func (g *Game) Token(pieceType uint8) shared.Token {
	return g.entities[pieceType]
}

func (g *Game) SheetImageSize() int {
	return g.sheetImageSize
}

func (g *Game) TranslateRFtoXY(rank, file uint8) (float64, float64) {
	return float64(uint(rank-1) * g.squareSize), float64(uint(8-file) * g.squareSize)
}

func (g *Game) TranslateXYtoRF(x, y int) (uint8, uint8, bool) {
	rank := uint8(float32(x-1)/float32(g.squareSize)) + 1
	file := 8 - uint8(float32(y-2)/float32(g.squareSize))
	if rank < 1 || rank > 8 || file < 1 || file > 8 {
		return 0, 0, false
	}
	return rank, file, true
}

func (g *Game) TranslateRFtoIndex(rank, file uint8) uint8 {
	index := (8-file)*8 + rank - 1
	return index
}

func (g *Game) TranslateIndexToRF(index uint8) (uint8, uint8) {
	rank := index%8 + 1
	file := 8 - index/8
	return rank, file
}

func (g *Game) TranslateIndexToXY(index uint8) (float64, float64) {
	rank, file := g.TranslateIndexToRF(index)
	x, y := g.TranslateRFtoXY(rank, file)
	return x, y
}

func (g *Game) TranslateRFtoN(rank, file uint8) string {
	if rank < 1 || rank > 8 || file < 1 || file > 8 {
		return ""
	}
	return string(notation[file-1]) + string('1'+rank-1)
}

func (g *Game) TranslateNtoRF(n string) (uint8, uint8, bool) {
	if len(n) != 2 {
		return 0, 0, false
	}
	n = strings.ToLower(n)
	file := n[0] - 'a' + 1
	rank := n[1] - '1' + 1
	if file < 1 || file > 8 || rank < 1 || rank > 8 {
		return 0, 0, false
	}
	return rank, file, true
}

func (g *Game) TranslateNtoIndex(n string) (uint8, bool) {
	if rank, file, ok := g.TranslateNtoRF(n); ok {
		return g.TranslateRFtoIndex(rank, file), true
	}
	return 0xFF, false
}

// -- DEBUG -------------------------------------------------------------

func (g *Game) EnableDebug() bool {
	return g.enabledDebug
}

func (g *Game) DebugX(rank uint8) int {
	return g.debugX[rank]
}

func (g *Game) DebugY() int {
	return g.debugY
}
func (g *Game) DebugFen() int {
	return g.fenY
}
