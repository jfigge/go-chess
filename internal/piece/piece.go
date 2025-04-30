package piece

import (
	"github.com/hajimehoshi/ebiten/v2"
	"us.figge.chess/internal/shared"
)

const (
	Pawn   uint = 0b00000001
	Knight uint = 0b00000010
	Bishop uint = 0b00000011
	Rook   uint = 0b00000100
	Queen  uint = 0b00000101
	King   uint = 0b00000110
)

type Piece struct {
	c        shared.Configuration
	token    shared.Token
	rank     int
	file     int
	op       *ebiten.DrawImageOptions
	captured bool
	promoted bool
	moved    bool
	inCheck  bool
}

func NewPiece(c shared.Configuration, pieceType uint, rank, file int) *Piece {
	piece := &Piece{
		c:        c,
		token:    c.Token(pieceType),
		rank:     rank,
		file:     file,
		captured: false,
		promoted: false,
		moved:    false,
		inCheck:  false,
	}
	piece.Position(rank, file)
	return piece
}

func (p *Piece) Draw(dst *ebiten.Image) {
	p.token.Draw(dst, p.op)
}

func (p *Piece) Position(rank, file int) {
	p.rank = rank
	p.file = file
	p.op = &ebiten.DrawImageOptions{}
	x, y := p.c.Translate(rank, file)
	p.op.GeoM.Translate(x, y)
}
