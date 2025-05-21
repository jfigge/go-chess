package engine

import (
	"fmt"
	"log"
	"strings"
	. "us.figge.chess/internal/common"
	"us.figge.chess/internal/engine/uci"
)

type Engine struct {
	position  *Position
	stockfish *uci.Engine
	fen       string
	moves     string
	cpuPlayer bool
}

func NewEngine() *Engine {
	e := &Engine{
		position: NewPosition(),
	}
	e.cpuPlayer = true

	var err error
	path := "/Users/jason/src/shell/commands/stockfish"
	e.stockfish, err = uci.NewEngine(path)
	if err != nil {
		log.Fatalf("Error launching engine: %v\n", err)
	}
	err = e.stockfish.UCI()
	if err != nil {
		log.Fatalln("Error setting UCI mode:", err)
	}
	engineOptions := uci.Options{
		MultiPV: 1,
		Hash:    1024,
		Ponder:  false,
		OwnBook: true,
		Threads: 6,
	}
	err = e.stockfish.SetOptions(engineOptions)
	if err != nil {
		log.Fatalln("Error setting engine options:", err)
	}
	return e
}

func (e *Engine) SetFEN(fen string) {
	fen = strings.TrimSpace(fen)
	if fen == "" {
		fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	}
	e.fen = fen
	e.position.SetupBoard(fen)
	err := e.stockfish.SetFEN(fen)
	if err != nil {
		log.Fatalf("Error setting FEN [%s]: %v\n", fen, err)
	}
}

func (e *Engine) GetBoards() []uint64 {
	boards := make([]uint64, 8, 8)
	for i := range 8 {
		boards[i] = e.position.bitboards[i]
	}
	return boards
}

func (e *Engine) GetEnPassant() (uint8, bool) {
	ep := e.position.EnPassant()
	if ep == 0 {
		return 0, false
	}
	return BtoI(ep), true
}

func (e *Engine) GetPieceType(rank, file uint8) (uint8, bool) {
	return e.position.identifyPiece(RFtoB(rank, file))
}

func (e *Engine) Turn() uint8 {
	return e.position.Turn()
}
func (e *Engine) MovePiece(from, to, pieceType uint8) (string, bool) {
	return e.position.MovePiece(from, to, pieceType)
	//if e.cpuPlayer {
	//	return true
	//}
	//e.moves += " " + RFtoN(ItoRF(from)) + RFtoN(ItoRF(to))
	//err := e.stockfish.SetMoves(e.moves)
	//if err != nil {
	//	log.Printf("Error setting move [%s]: %v\n", e.moves, err)
	//}
	//var results *uci.Results
	//results, err = e.stockfish.Go(10, "", int64(time.Second))
	//if err != nil {
	//	log.Printf("Error getting moves: %v\n", err)
	//}
	//fmt.Printf("Best move: %s\t Mate: %v\n", results.BestMove, results.Results[0].Mate)
	//rank, file, _ := NtoRF(results.BestMove[:2])
	//from = RFtoI(rank, file)
	//rank, file, _ = NtoRF(results.BestMove[2:4])
	//to = RFtoI(rank, file)
	//pieceType, _ = e.position.identifyPiece(ItoB(from))
	//msg, ok = e.position.MovePiece(from, to, pieceType)
	//fmt.Println(msg)
	//if !ok {
	//	log.Printf("Error moving piece from %d to %d\n", from, to)
	//	return false
	//}
	//e.moves += " " + results.BestMove
	//return true
}

func (e *Engine) showPieces(pieceType uint8) {
	fmt.Println("Board")
	e.position.debugPrintBoard()
	bb, pb := PTtoBB(pieceType)
	piece := ""
	switch bb - 2 {
	case PiecePawn:
		piece = "Pawn"
	case PieceKnight:
		piece = "Knight"
	case PieceBishop:
		piece = "Bishop"
	case PieceRook:
		piece = "Rook"
	case PieceQueen:
		piece = "Queen"
	case PieceKing:
		piece = "King"
	}
	fmt.Println("Piece: " + piece)
	e.position.debugPrintBitBoard(e.position.bitboards[bb], 0)
	player := "White"
	if pb == PlayerBlack {
		player = "Black"
	}
	fmt.Println("Player: ", player)
	e.position.debugPrintBitBoard(e.position.bitboards[pb], 0)
}

func (e *Engine) GetMoves(rank uint8, file uint8, pieceType uint8) []uint8 {
	return []uint8{44, 36, 28}
}
