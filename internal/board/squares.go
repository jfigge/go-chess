package board

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"us.figge.chess/internal/piece"
)

type square struct {
	piece     *piece.Piece
	size      float32
	x         float32
	y         float32
	index     int
	highlight bool
	valid     bool
	invalid   bool
}

type squares struct {
	squares [64]square
}

func (s *square) SetPiece(p *piece.Piece) {
	s.piece = p
}

func (s *square) Draw(b *Board, target *ebiten.Image) {
	c := b.ColorHighlight()
	if s.valid {
		c = b.ColorValid()
	} else if s.invalid {
		c = b.ColorInvalid()
	}
	if c != nil {
		vector.DrawFilledRect(target, s.x, s.y, s.size, s.size, c, false)
	}
}
