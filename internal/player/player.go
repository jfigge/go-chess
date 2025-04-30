package player

import (
	"github.com/hajimehoshi/ebiten/v2"
	"us.figge.chess/internal/piece"
	. "us.figge.chess/internal/shared"
)

type Player struct {
	c      Configuration
	color  uint8
	pieces []*piece.Piece
	KSC    bool
	QSC    bool
}

func NewPlayer(c Configuration, color uint8) *Player {
	player := &Player{
		c:      c,
		color:  color,
		pieces: []*piece.Piece{},
		KSC:    true,
		QSC:    true,
	}
	return player
}

func (p *Player) Draw(target *ebiten.Image) {
	for i := 0; i < len(p.pieces); i++ {
		p.pieces[i].Draw(target)
	}
}

func (p *Player) AddPiece(pieceType uint8, rank, file uint8) {
	newPiece := piece.NewPiece(p.c, pieceType|p.color)
	newPiece.Position(rank, file)
	p.pieces = append(p.pieces, newPiece)
}

func (p *Player) SetKingsideCastle(ksc bool) {
	p.KSC = ksc
}
func (p *Player) SetQueensideCastle(qsc bool) {
	p.QSC = qsc
}
