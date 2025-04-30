package player

import (
	"github.com/hajimehoshi/ebiten/v2"
	"us.figge.chess/internal/piece"
	"us.figge.chess/internal/shared"
)

const (
	White uint = 0b00001000
	Black uint = 0b00010000
)

type Player struct {
	c      shared.Configuration
	color  uint
	pieces [16]*piece.Piece
}

func NewPlayer(c shared.Configuration, color uint) *Player {
	pawnRank := 2
	backRank := 1
	if color == Black {
		pawnRank = 7
		backRank = 8
	}
	player := &Player{
		c:     c,
		color: uint(color),
		pieces: [16]*piece.Piece{
			piece.NewPiece(c, piece.Pawn|color, pawnRank, 1),
			piece.NewPiece(c, piece.Pawn|color, pawnRank, 2),
			piece.NewPiece(c, piece.Pawn|color, pawnRank, 3),
			piece.NewPiece(c, piece.Pawn|color, pawnRank, 4),
			piece.NewPiece(c, piece.Pawn|color, pawnRank, 5),
			piece.NewPiece(c, piece.Pawn|color, pawnRank, 6),
			piece.NewPiece(c, piece.Pawn|color, pawnRank, 7),
			piece.NewPiece(c, piece.Pawn|color, pawnRank, 8),
			piece.NewPiece(c, piece.Rook|color, backRank, 1),
			piece.NewPiece(c, piece.Knight|color, backRank, 2),
			piece.NewPiece(c, piece.Bishop|color, backRank, 3),
			piece.NewPiece(c, piece.Queen|color, backRank, 4),
			piece.NewPiece(c, piece.King|color, backRank, 5),
			piece.NewPiece(c, piece.Bishop|color, backRank, 6),
			piece.NewPiece(c, piece.Knight|color, backRank, 7),
			piece.NewPiece(c, piece.Rook|color, backRank, 8),
		},
	}
	return player
}

func (p *Player) Draw(screen *ebiten.Image) {
	for i := 0; i < 16; i++ {
		p.pieces[i].Draw(screen)
	}
}
