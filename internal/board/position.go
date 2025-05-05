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

func (p *position) turn() uint8 {
	return uint8(p.bitboards[common.Pawn|common.White]) & common.TurnMask
}
func (p *position) castleRights() uint8 {
	return uint8(p.bitboards[common.Pawn|common.White]) & common.CastleRightsMask
}
func (p *position) enPassant() uint8 {
	return uint8(p.bitboards[common.Pawn|common.Black] >> 56)
}
func (p *position) setTurn(turn uint8) {
	p.bitboards[common.Pawn|common.White] &= ^uint64(common.TurnMask)
	p.bitboards[common.Pawn|common.White] |= uint64(turn)
}
func (p *position) setCastleRights(castleRights uint8) {
	p.bitboards[common.Pawn|common.White] &= ^uint64(common.CastleRightsMask)
	p.bitboards[common.Pawn|common.White] |= uint64(castleRights)
}
func (p *position) setEnPassant(enPassant uint8) {
	p.bitboards[common.Pawn|common.Black] &= ^uint64(common.EnPassantMask)
	p.bitboards[common.Pawn|common.Black] |= uint64(enPassant) << 56
}

func (p *position) setPiece(pieceType uint8, rank, file int) {
	index := p.TranslateRFtoIndex(rank, file)
	if index < 0 {
		panic("Invalid bit position")
	}
	bitboards[pieceType] |= uint64(1) << index
}
func (p *position) removePiece(pieceType uint8, rank, file int) {
	index := p.TranslateRFtoIndex(rank, file)
	if index < 0 {
		panic("Invalid bit position")
	}
	bitboards[pieceType] &= ^(uint64(1) << index)
}
func (p *position) clearPiece(rank, file int) {
	index := p.TranslateRFtoIndex(rank, file)
	if index < 0 {
		panic("Invalid bit position")
	}
	if b, ok := p.findPiece(index); ok {
		bitboards[b] &= ^(uint64(1) << index)
	}
}
func (p *position) findPiece(index int) (uint8, bool) {
	for b := range uint8(12) {
		if bitboards[b]&(1<<index) != 0 {
			return b, true
		}
	}
	return 0, false
}

func (p *position) SetupBoard(fen string) {
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
			p.setPiece(fenPieceMap[c], rank, file)
			index++
		}
	}
	parts := strings.Split(fen, " ")
	if len(parts) > 0 && len(parts[0]) > 0 {
		p.readFenTurn(parts[0])
	}
	if len(parts) > 1 && len(parts[1]) > 0 {
		p.readFenCastling(parts[1])
	}
	if len(parts) > 2 && len(parts[2]) > 0 {
		p.readFenEnpassant(parts[2])
	}
	if len(parts) > 3 && len(parts[3]) > 0 {
		p.readFenHalfMove(parts[3])
	}
	if len(parts) > 4 && len(parts[4]) > 0 {
		p.readFenFullMove(parts[4])
	}
}
func (p *position) readFenTurn(turn string) {
	t := common.White
	if strings.EqualFold(turn, "b") {
		t = common.Black
	} else if turn != "w" {
		fmt.Printf("Invalid turn fen: %s\n", turn)
	}
	p.setTurn(t)
}
func (p *position) readFenCastling(castling string) {
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
	p.setCastleRights(castleRights)
}
func (p *position) readFenEnpassant(enpassant string) {
	switch enpassant {
	case "-":
		p.setEnPassant(uint8(0))
	default:
		if _, file, ok := p.TranslateNtoRF(enpassant); ok {
			p.setEnPassant(uint8(9 - file))
		} else {
			fmt.Printf("Invalid enPassant fen: %s\n", enpassant)
		}
	}
}
func (p *position) readFenHalfMove(halfMove string) {
	halfMoveCount := 0
	if _, err := fmt.Sscanf(halfMove, "%d", &halfMoveCount); err != nil {
		fmt.Printf("Invalid halfmove fen: %s\n", halfMove)
	}
	p.halfmoves = make([]uint64, halfMoveCount)
}
func (p *position) readFenFullMove(fullMove string) {
	fullMoveCount := 0
	if _, err := fmt.Sscanf(fullMove, "%d", &fullMoveCount); err != nil {
		fmt.Printf("Invalid fullmove fen: %s\n", fullMove)
	}
	p.fullmoves = fullMoveCount
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
func (p *position) RecordBoardFen() string {
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
	p.writeFenTurn(&sb)
	p.writeFenCastling(&sb)
	p.writeFenEnpassant(&sb)
	p.writeFenHalfMove(&sb)
	p.writeFenFullMove(&sb)
	return sb.String()
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
func (p *position) writeFenTurn(sb *strings.Builder) {
	if p.turn() == common.White {
		sb.WriteString(" w")
	} else {
		sb.WriteString(" b")
	}
}
func (p *position) writeFenCastling(sb *strings.Builder) {
	sb.WriteByte(' ')
	castlingRights := p.castleRights()
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
func (p *position) writeFenEnpassant(sb *strings.Builder) {
	ep := p.enPassant()
	if ep == 0 {
		sb.WriteString(" -")
	} else {
		notation := string([]byte{'a' + 8 - ep})
		if p.turn() == common.White {
			notation += "6"
		} else {
			notation += "3"
		}
		sb.WriteString(" " + notation)
	}
}
func (p *position) writeFenHalfMove(sb *strings.Builder) {
	sb.WriteByte(' ')
	sb.WriteString(fmt.Sprintf("%d", len(p.halfmoves)))
}
func (p *position) writeFenFullMove(sb *strings.Builder) {
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
