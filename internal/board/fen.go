package board

import (
	"fmt"
	"strings"
	"us.figge.chess/internal/piece"
	"us.figge.chess/internal/player"
	. "us.figge.chess/internal/shared"
)

func (b *Board) Fen() string {
	sb := strings.Builder{}

	index := 0
	count := 0
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			p := b.squares[index].piece
			if p != nil {
				b.WriteFenEntry(&sb, &count, p)
			} else {
				count++
			}
			index++
		}
		if count > 0 {
			b.WriteFenEntry(&sb, &count, nil)
		}
		if i < 7 {
			sb.WriteByte('/')
		}
	}
	b.WriteTurn(&sb)
	b.WriteCastling(&sb)
	b.WriteEnpassant(&sb)
	b.WriteHalfMove(&sb)
	b.WriteFullMove(&sb)
	fmt.Println(sb.String())
	return sb.String()
}
func (b *Board) WriteFenEntry(sb *strings.Builder, count *int, p *piece.Piece) {
	if *count > 0 {
		sb.WriteByte(byte(*count + '0'))
	}
	if p != nil {
		sb.WriteByte(p.Fen())
	}
	*count = 0
}
func (b *Board) WriteTurn(sb *strings.Builder) {
	if b.turn == White {
		sb.WriteString(" w")
	} else {
		sb.WriteString(" b")
	}
}
func (b *Board) WriteCastling(sb *strings.Builder) {
	sb.WriteByte(' ')
	castlingAvailable := false
	if b.players[0].KingsideCastle() {
		sb.WriteByte('K')
		castlingAvailable = true
	}
	if b.players[0].QueensideCastle() {
		sb.WriteByte('Q')
		castlingAvailable = true
	}
	if b.players[1].KingsideCastle() {
		sb.WriteByte('k')
		castlingAvailable = true
	}
	if b.players[1].QueensideCastle() {
		sb.WriteByte('q')
		castlingAvailable = true
	}
	if !castlingAvailable {
		sb.WriteByte('-')
	}
}
func (b *Board) WriteEnpassant(sb *strings.Builder) {
	if b.enpassant == 0xff {
		sb.WriteString(" -")
	} else {
		rank, file := b.TranslateIndexToRF(b.enpassant)
		sb.WriteByte(' ')
		sb.WriteString(b.TranslateRFtoN(rank, file))
	}
}
func (b *Board) WriteHalfMove(sb *strings.Builder) {
	sb.WriteByte(' ')
	sb.WriteString(fmt.Sprintf("%d", b.halfMove))
}
func (b *Board) WriteFullMove(sb *strings.Builder) {
	sb.WriteByte(' ')
	sb.WriteString(fmt.Sprintf("%d", b.fullMove))
}

func (b *Board) SetFen(fen string) {
	b.resetBoard()
	b.fen = fen
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
				rank, file := b.TranslateIndexToRF(index)
				b.squares[index].piece = b.players[p&0b1].AddPiece(p&0b1110, rank, file)
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
		if e, ok := b.TranslateNtoIndex(enpassant); ok {
			b.enpassant = e
		} else {
			fmt.Printf("Invalid enpassant fen: %s\n", enpassant)
		}
	}
}
func (b *Board) setHalfMove(halfMove string) {
	halfMoveCount := 0
	if _, err := fmt.Sscanf(halfMove, "%d", &halfMoveCount); err != nil {
		fmt.Printf("Invalid halfmove fen: %s\n", halfMove)
	}
	b.halfMove = halfMoveCount
}
func (b *Board) setFullMove(fullMove string) {
	fullMoveCount := 0
	if _, err := fmt.Sscanf(fullMove, "%d", &fullMoveCount); err != nil {
		fmt.Printf("Invalid fullmove fen: %s\n", fullMove)
	}
	b.fullMove = fullMoveCount
}
