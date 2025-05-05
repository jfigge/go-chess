package board

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"strings"
	"us.figge.chess/internal/common"
)

var (
	keys        [768]uint64
	bitboards   [12]uint64
	fenPieceMap = map[byte]uint8{
		'P': common.White | common.Pawn,
		'N': common.White | common.Knight,
		'B': common.White | common.Bishop,
		'R': common.White | common.Rook,
		'Q': common.White | common.Queen,
		'K': common.White | common.King,
		'p': common.Black | common.Pawn,
		'n': common.Black | common.Knight,
		'b': common.Black | common.Bishop,
		'r': common.Black | common.Rook,
		'q': common.Black | common.Queen,
		'k': common.Black | common.King,
	}
)

type position struct {
	common.Configuration
	hashKey   string
	halfmoves []uint64
	fullmoves int
	bitboards [12]uint64
}

func (p *position) Turn() uint8 {
	return uint8(p.bitboards[common.Pawn|common.White]) & common.TurnMask
}
func (p *position) CastleRights() uint8 {
	return uint8(p.bitboards[common.Pawn|common.White]) & common.CastleRightsMask
}
func (p *position) EnPassant() uint8 {
	return uint8(p.bitboards[common.Pawn|common.Black] >> 56)
}
func (p *position) SetTurn(turn uint8) {
	p.bitboards[common.Pawn|common.White] &= ^uint64(common.TurnMask)
	p.bitboards[common.Pawn|common.White] |= uint64(turn)
}
func (p *position) SetCastleRights(castleRights uint8) {
	p.bitboards[common.Pawn|common.White] &= ^uint64(common.CastleRightsMask)
	p.bitboards[common.Pawn|common.White] |= uint64(castleRights)
}
func (p *position) SetEnPassant(enPassant uint8) {
	p.bitboards[common.Pawn|common.Black] &= ^uint64(common.EnPassantMask)
	p.bitboards[common.Pawn|common.Black] |= uint64(enPassant) << 56
}
func (p *position) SetPiece(pieceType uint8, rank, file int) {
	index := p.TranslateRFtoIndex(rank, file)
	if index < 0 {
		panic("Invalid bit position")
	}
	bitboards[pieceType] |= uint64(1) << index
}
func (p *position) RemovePiece(pieceType uint8, rank, file int) {
	index := p.TranslateRFtoIndex(rank, file)
	if index < 0 {
		panic("Invalid bit position")
	}
	bitboards[pieceType] &= ^(uint64(1) << index)
}
func (p *position) ClearSquare(rank, file int) {
	index := p.TranslateRFtoIndex(rank, file)
	if index < 0 {
		panic("Invalid bit position")
	}
	if b, ok := p.findPiece(index); ok {
		bitboards[b] &= ^(uint64(1) << index)
	}
}

func (p *position) Hash() uint64 {
	var hash uint64
	for i := range bitboards {
		board := bitboards[i]
		index := 0
		for board != 0 {
			if board&1 != 0 {
				hash ^= keys[index]
			}
			board >>= 1
			index++
		}
	}
	return hash
}

func (p *position) setupBoard(fen string) {
	fen = strings.TrimSpace(fen)
	if fen == "" {
		fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq _ 0 1"
	}
	p.fullmoves = 0
	p.halfmoves = make([]uint64, 0)
	p.bitboards = [12]uint64{}

	// Parse the FEN string and set up the board
	index := 0
bitboards:
	for i := 0; i < len(fen); i++ {
		c := fen[i]
		switch c {
		case '1', '2', '3', '4', '5', '6', '7', '8':
			index += int(c - '0')
		case '/':
			continue
		case ' ':
			fen = fen[i+1:]
			break bitboards
		case 'K', 'Q', 'R', 'B', 'N', 'P', 'k', 'q', 'r', 'b', 'n', 'p':
			rank, file := p.TranslateIndexToRF(index)
			p.SetPiece(fenPieceMap[c], rank, file)
			index++
		}
	}
	parts := strings.Split(fen, " ")
	if len(parts) > 0 && len(parts[0]) > 0 {
		p.setTurn(parts[0])
	}
	if len(parts) > 1 && len(parts[1]) > 0 {
		p.setCastling(parts[1])
	}
	if len(parts) > 2 && len(parts[2]) > 0 {
		p.setEnpassant(parts[2])
	}
	if len(parts) > 3 && len(parts[3]) > 0 {
		p.setHalfMove(parts[3])
	}
	if len(parts) > 4 && len(parts[4]) > 0 {
		p.setFullMove(parts[4])
	}
}
func (p *position) setTurn(turn string) {
	t := common.White
	if strings.EqualFold(turn, "b") {
		t = common.Black
	} else if turn != "w" {
		fmt.Printf("Invalid turn fen: %s\n", turn)
	}
	p.SetTurn(t)
}
func (p *position) setCastling(castling string) {
	castleRights := uint8(0)
	for _, c := range castling {
		switch c {
		case 'K':
			castleRights |= common.CastleRightsWhiteKing
		case 'Q':
			castleRights |= common.CastleRightsWhiteQueen
		case 'k':
			castleRights |= common.CastleRightsBlackKing
		case 'q':
			castleRights |= common.CastleRightsBlackQueen
		case '-':
			castleRights = 0
			break
		default:
			fmt.Printf("Invalid castling fen: %s\n", castling)
		}
	}
	p.SetCastleRights(castleRights)
}
func (p *position) setEnpassant(enpassant string) {
	switch enpassant {
	case "-":
		p.SetEnPassant(uint8(0))
	default:
		if _, file, ok := p.TranslateNtoRF(enpassant); ok {
			p.SetEnPassant(uint8(9 - file))
		} else {
			fmt.Printf("Invalid enpassant fen: %s\n", enpassant)
		}
	}
}
func (p *position) setHalfMove(halfMove string) {
	halfMoveCount := 0
	if _, err := fmt.Sscanf(halfMove, "%d", &halfMoveCount); err != nil {
		fmt.Printf("Invalid halfmove fen: %s\n", halfMove)
	}
	p.halfmoves = make([]uint64, halfMoveCount)
}
func (p *position) setFullMove(fullMove string) {
	fullMoveCount := 0
	if _, err := fmt.Sscanf(fullMove, "%d", &fullMoveCount); err != nil {
		fmt.Printf("Invalid fullmove fen: %s\n", fullMove)
	}
	p.fullmoves = fullMoveCount
}

func (p *position) generateFen() string {
	sb := strings.Builder{}
	count := 0
	for rank := 1; rank <= 8; rank++ {
		for file := 1; file <= 8; file++ {
			index := p.TranslateRFtoIndex(9-rank, file)
			if pieceType, ok := p.findPiece(index); ok {
				p.writeFenEntry(&sb, &count, p.Piece(pieceType))
			} else {
				count++
			}
		}
		if count != 0 {
			p.writeFenEntry(&sb, &count, nil)
		}
		if rank < 8 {
			sb.WriteByte('/')
		}
	}
	p.writeTurn(&sb)
	p.writeCastling(&sb)
	p.writeEnpassant(&sb)
	p.writeHalfMove(&sb)
	p.writeFullMove(&sb)
	return sb.String()
}
func (p *position) findPiece(index int) (uint8, bool) {
	for b := range uint8(12) {
		if bitboards[b]&(1<<index) != 0 {
			return b, true
		}
	}
	return 0, false
}
func (p *position) writeFenEntry(sb *strings.Builder, count *int, piece common.Piece) {
	if *count > 0 {
		sb.WriteByte(byte(*count + '0'))
	}
	if piece != nil {
		sb.WriteByte(piece.Fen())
	}
	*count = 0
}
func (p *position) writeTurn(sb *strings.Builder) {
	if p.Turn() == common.White {
		sb.WriteString(" w")
	} else {
		sb.WriteString(" b")
	}
}
func (p *position) writeCastling(sb *strings.Builder) {
	sb.WriteByte(' ')
	castlingRights := p.CastleRights()
	if castlingRights == 0 {
		sb.WriteByte('-')
		return
	}
	if castlingRights&common.CastleRightsWhiteKing > 0 {
		sb.WriteByte('K')
	}
	if castlingRights&common.CastleRightsWhiteQueen > 0 {
		sb.WriteByte('Q')
	}
	if castlingRights&common.CastleRightsBlackKing > 0 {
		sb.WriteByte('k')
	}
	if castlingRights&common.CastleRightsBlackQueen > 0 {
		sb.WriteByte('q')
	}
}
func (p *position) writeEnpassant(sb *strings.Builder) {
	enPassant := p.EnPassant()
	if enPassant == 0 {
		sb.WriteString(" -")
	} else {
		notation := string([]byte{'a' + 8 - enPassant})
		if p.Turn() == common.White {
			notation += "6"
		} else {
			notation += "3"
		}
		sb.WriteString(" " + notation)
	}
}
func (p *position) writeHalfMove(sb *strings.Builder) {
	sb.WriteByte(' ')
	sb.WriteString(fmt.Sprintf("%d", len(p.halfmoves)))
}
func (p *position) writeFullMove(sb *strings.Builder) {
	sb.WriteByte(' ')
	sb.WriteString(fmt.Sprintf("%d", p.fullmoves))
}

func (p *position) drawPieces(dst *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	for rank := 1; rank <= 8; rank++ {
		for file := 1; file <= 8; file++ {
			index := p.TranslateRFtoIndex(rank, file)
			if b, ok := p.findPiece(index); ok {
				op.GeoM.Reset()
				x, y := p.TranslateRFtoXY(rank, file)
				op.GeoM.Translate(x, y)
				p.Piece(b).Draw(dst, op)
			}
		}
		fmt.Println()
	}
}
func (p *position) debugPrintBoard() {
	for rank := 1; rank <= 8; rank++ {
		fmt.Printf("%d  ", 9-rank)
		for file := 1; file <= 8; file++ {
			index := p.TranslateRFtoIndex(9-rank, file)
			str := ". "
			for b := range uint8(12) {
				if bitboards[b]&(1<<index) != 0 {
					str = fmt.Sprintf("%s ", string([]byte{p.Piece(b).Fen()}))
					break
				}
			}
			fmt.Print(str)
		}
		fmt.Println()
	}
	fmt.Print("\n   a b c d e f g h\n")
}

func init() {
	// Key generation
	bs := make([]byte, 8)
	for i := range keys {
		n, err := rand.Read(bs)
		if err != nil || n != 8 {
			panic("Random key generation failed")
		}
		keys[i] = binary.BigEndian.Uint64(bs)
	}
}
