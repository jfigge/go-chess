package engine

import (
	"fmt"
	. "us.figge.chess/internal/common"
)

type Engine struct {
	position *Position
}

func NewEngine() *Engine {
	e := &Engine{
		position: NewPosition(),
	}
	return e
}

func (e *Engine) SetFEN(fen string) {
	e.position.SetupBoard(fen)
}

func (e *Engine) GetBoards() []uint64 {
	boards := make([]uint64, 8, 8)
	for i := range 8 {
		boards[i] = e.position.bitboards[i]
	}
	return boards
}

func (e *Engine) GetPieceType(rank, file uint8) (uint8, bool) {
	return e.position.identifyPiece(RFtoB(rank, file))
}

func (e *Engine) MovePiece(from, to, pieceType uint8) {
	e.position.MovePiece(from, to, pieceType)
	fmt.Printf("Move: %s\n", RFtoN(ItoRF(to)))
}
