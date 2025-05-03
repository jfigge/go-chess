package game

import (
	"strings"
	"us.figge.chess/internal/common"
)

const (
	notation = "abcdefgh"
)

func (g *Game) SquareSize() int {
	return g.squareSize
}

func (g *Game) Token(pieceType uint8) common.Token {
	return g.entities[pieceType]
}

func (g *Game) SheetImageSize() int {
	return g.sheetImageSize
}

func (g *Game) Turn(colorType uint8) string {
	return colors[colorType]
}

func (g *Game) HighlightAttacks() bool {
	return g.highlightAttacks
}

func (g *Game) ShowStrength() bool {
	return g.showStrength
}

func (g *Game) FontHeight() int {
	return g.fontHeight
}

func (g *Game) TranslateRFtoXY(rank, file int) (float64, float64) {
	return float64((file - 1) * g.squareSize), float64((8 - rank) * g.squareSize)
}

func (g *Game) TranslateXYtoRF(x, y int) (int, int, bool) {
	if x < 0 || y < 0 {
		return 0, 0, false
	}
	rank := 7 - int(float32(y-2)/float32(g.squareSize))
	file := int(float32(x-1) / float32(g.squareSize))
	if rank < 0 || rank > 7 || file < 0 || file > 7 {
		return 0, 0, false
	}
	return rank + 1, file + 1, true
}

func (g *Game) TranslateRFtoIndex(rank, file int) int {
	index := (8-rank)*8 + file - 1
	return index
}

func (g *Game) TranslateIndexToRF(index int) (int, int) {
	rank := 8 - index/8
	file := index%8 + 1
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

func (g *Game) DebugX(file int) int {
	return g.debugX[file]
}

func (g *Game) DebugY() int {
	return g.debugY
}
