package player

import (
	"github.com/hajimehoshi/ebiten/v2"
	. "us.figge.chess/internal/common"
	"us.figge.chess/internal/piece"
)

type Player struct {
	Configuration
	color  uint8
	pieces []*piece.Piece
	ksc    bool
	qsc    bool
}

func NewPlayer(c Configuration, color uint8) *Player {
	player := &Player{
		Configuration: c,
		color:         color,
		pieces:        []*piece.Piece{},
		ksc:           true,
		qsc:           true,
	}
	return player
}

func (p *Player) Draw(target *ebiten.Image) {
	for _, e := range p.pieces {
		e.Draw(target, true)
	}
}

func (p *Player) AddPiece(pieceType uint8, rank, file int) *piece.Piece {
	newPiece := piece.NewPiece(Configuration(p), pieceType|p.color)
	newPiece.Position(rank, file)
	p.pieces = append(p.pieces, newPiece)
	return newPiece
}

func (p *Player) SetKingsideCastle(ksc bool) {
	p.ksc = ksc
}
func (p *Player) SetQueensideCastle(qsc bool) {
	p.qsc = qsc
}
func (p *Player) KingsideCastle() bool {
	return p.ksc
}
func (p *Player) QueensideCastle() bool {
	return p.qsc
}
