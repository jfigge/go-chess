package game

import (
	"strings"
	"us.figge.chess/internal/shared"
)

const (
	notation = "abcdefgh"
)

func (g *Game) SquareSize() int {
	return g.squareSize
}

func (g *Game) Token(pieceType uint8) shared.Token {
	return g.entities[pieceType]
}

func (g *Game) SheetImageSize() int {
	return g.sheetImageSize
}

func (g *Game) TranslateRFtoXY(rank, file int) (float64, float64) {
	return float64((rank - 1) * g.squareSize), float64((8 - file) * g.squareSize)
}

func (g *Game) TranslateXYtoRF(x, y int) (int, int, bool) {
	rank := int(float32(x-1)/float32(g.squareSize)) + 1
	file := 8 - int(float32(y-2)/float32(g.squareSize))
	if rank < 1 || rank > 8 || file < 1 || file > 8 {
		return 0, 0, false
	}
	return rank, file, true
}

func (g *Game) TranslateRFtoIndex(rank, file int) int {
	index := (8-file)*8 + rank - 1
	return index
}

func (g *Game) TranslateIndexToRF(index int) (int, int) {
	rank := index%8 + 1
	file := 8 - index/8
	return rank, file
}

func (g *Game) TranslateIndexToXY(index int) (float64, float64) {
	rank, file := g.TranslateIndexToRF(index)
	x, y := g.TranslateRFtoXY(rank, file)
	return x, y
}

func (g *Game) TranslateRFtoN(rank, file int) string {
	if rank < 1 || rank > 8 || file < 1 || file > 8 {
		return ""
	}
	return string(notation[file-1]) + string('1'+uint8(rank-1))
}

func (g *Game) TranslateNtoRF(n string) (int, int, bool) {
	if len(n) != 2 {
		return 0, 0, false
	}
	n = strings.ToLower(n)
	file := int(n[0] - 'a' + 1)
	rank := int(n[1] - '1' + 1)
	if file < 1 || file > 8 || rank < 1 || rank > 8 {
		return 0, 0, false
	}
	return rank, file, true
}

func (g *Game) TranslateNtoIndex(n string) (int, bool) {
	if rank, file, ok := g.TranslateNtoRF(n); ok {
		return g.TranslateRFtoIndex(rank, file), true
	}
	return 0xFF, false
}

// -- DEBUG -------------------------------------------------------------

func (g *Game) EnableDebug() bool {
	return g.debugEnabled
}

func (g *Game) DebugX(rank int) int {
	return g.debugX[rank]
}

func (g *Game) DebugY() int {
	return g.debugY
}
func (g *Game) DebugFen() int {
	return g.fenY
}
