package piece

import (
	"github.com/hajimehoshi/ebiten/v2"
	. "us.figge.chess/internal/shared"
)

type Piece struct {
	Configuration
	Token
	rank uint8
	file uint8
	op   *ebiten.DrawImageOptions
}

func NewPiece(c Configuration, pieceType uint8) *Piece {
	piece := &Piece{
		Configuration: c,
		Token:         c.Token(pieceType),
		rank:          0,
		file:          0,
	}
	return piece
}

func (p *Piece) Draw(target *ebiten.Image) {
	p.Token.Draw(target, p.op)
}

func (p *Piece) Rank() uint8 {
	return p.rank
}
func (p *Piece) File() uint8 {
	return p.file
}

func (p *Piece) Position(rank, file uint8) {
	p.rank = rank
	p.file = file
	p.op = &ebiten.DrawImageOptions{}
	x, y := p.Translate(rank, file)
	p.op.GeoM.Translate(x, y)
}
