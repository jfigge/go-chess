package board

import (
	"fmt"
	"strings"
	"us.figge.chess/internal/player"
	. "us.figge.chess/internal/shared"
)

func (b *Board) Fen() string {
	return ""
}

func (b *Board) SetFen(fen string) {
	b.fen = fen
	b.turn = White
	b.fullMove = 0
	b.halfMove = 0
	b.enpassant = 0xff
	b.dragIndex = 0xff
	b.dragPiece = nil
	b.players[0] = player.NewPlayer(Configuration(b), White)
	b.players[1] = player.NewPlayer(Configuration(b), Black)
	b.resetBoard()
	index := 0
	parts := strings.Split(strings.TrimSpace(fen), " ")
	if len(parts) > 0 && len(parts[0]) > 0 {
		for i := 0; i < len(parts[0]); i++ {
			c := parts[0][i]
			switch c {
			case '1', '2', '3', '4', '5', '6', '7', '8':
				index += int(c - '0')
			case '/':
				continue
			case ' ':
				break
			case 'K', 'Q', 'R', 'B', 'N', 'P', 'k', 'q', 'r', 'b', 'n', 'p':
				p := FenPieceMap[c]
				rank, file := b.TranslateIndexToRF(uint8(index))
				b.board[index].piece = b.players[(p&Black)>>4].AddPiece(p, rank, file)
				index++
			}
		}
	}
	if len(parts) > 1 && len(parts[1]) > 0 {
		b.setTurn(parts[1])
	}
	if len(parts) > 2 && len(parts[2]) > 0 {
		b.setCastling(parts[2], b.players[0], b.players[1])
	}
	if len(parts) > 3 && len(parts[3]) > 0 {
		b.setEnpassant(parts[3])
	}
	if len(parts) > 4 && len(parts[4]) > 0 {
		b.setHalfMove(parts[4])
	}
	if len(parts) > 5 && len(parts[5]) > 0 {
		b.setFullMove(parts[5])
	}
}

func (b *Board) setTurn(turn string) {
	b.turn = White
	if turn == "b" {
		b.turn = Black
	} else if turn != "w" {
		fmt.Printf("Invalid turn fen: %s\n", turn)
	}
}

func (b *Board) setCastling(castling string, white, black *player.Player) {
	wksc := false
	wqsc := false
	bksc := false
	bqsc := false
	for _, c := range castling {
		switch c {
		case 'K':
			wksc = true
		case 'Q':
			wqsc = true
		case 'k':
			bksc = true
		case 'q':
			bqsc = true
		case '-': //ignore
		default:
			fmt.Printf("Invalid castling fen: %s\n", castling)
		}
	}
	white.SetKingsideCastle(wksc)
	white.SetQueensideCastle(wqsc)
	black.SetKingsideCastle(bksc)
	black.SetQueensideCastle(bqsc)
}

func (b *Board) setEnpassant(enpassant string) {
	switch enpassant {
	case "-":
		b.enpassant = 0xff
	default:
		rank := uint(enpassant[1] - '0')
		file := uint(enpassant[0] - 'a')
		if rank < 1 || rank > 8 || file < 0 || file > 7 {
			fmt.Printf("Invalid enpassant fen: %s\n", enpassant)
			return
		}
		b.enpassant = uint8((rank-1)*8 + file)
	}
}

func (b *Board) setHalfMove(halfMove string) {
	halfMoveCount := uint8(0)
	if _, err := fmt.Sscanf(halfMove, "%d", &halfMoveCount); err != nil {
		fmt.Printf("Invalid halfmove fen: %s\n", halfMove)
	}
	b.halfMove = halfMoveCount
}

func (b *Board) setFullMove(fullMove string) {
	fullMoveCount := uint(0)
	if _, err := fmt.Sscanf(fullMove, "%d", &fullMoveCount); err != nil {
		fmt.Printf("Invalid fullmove fen: %s\n", fullMove)
	}
	b.fullMove = fullMoveCount
}
