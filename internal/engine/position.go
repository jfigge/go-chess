package engine

import (
	"fmt"
	"strings"
	. "us.figge.chess/internal/common"
)

var (
	//hashKey        string
	//keys           [576]uint64
	knightMoves    [64]uint64
	kingMoves      [64]uint64
	whitePawnMoves [64]uint64
	blackPawnMoves [64]uint64
	ranks          [8]uint64
	files          [8]uint64
	fenPieceMap    = map[byte]uint8{
		'P': PlayerWhite | PiecePawn,
		'N': PlayerWhite | PieceKnight,
		'B': PlayerWhite | PieceBishop,
		'R': PlayerWhite | PieceRook,
		'Q': PlayerWhite | PieceQueen,
		'K': PlayerWhite | PieceKing,
		'p': PlayerBlack | PiecePawn,
		'n': PlayerBlack | PieceKnight,
		'b': PlayerBlack | PieceBishop,
		'r': PlayerBlack | PieceRook,
		'q': PlayerBlack | PieceQueen,
		'k': PlayerBlack | PieceKing,
	}
)

const (
	fen       = "PpNnBbRrQqKk"
	algebraic = " NBRQK"
	RANK8     = 0xFF00000000000000
	RANK1     = 0x00000000000000FF
	FileA     = 0x0101010101010101
	FileH     = 0x8080808080808080
)

type Position struct {
	status    uint8
	bitboards [9]uint64
	halfMoves []uint64
	fullMoves int
}

func NewPosition() *Position {
	return &Position{
		status:    PlayerWhite | CastleRightsMask,
		bitboards: [9]uint64{},
		halfMoves: make([]uint64, 0),
		fullMoves: 0,
	}
}

func (p *Position) Pieces(turn uint8) uint64 {
	turn = turn & PlayerMask
	return p.bitboards[turn]
}

func (p *Position) Turn() uint8 {
	return p.status & PlayerMask
}
func (p *Position) CastleRights() uint8 {
	return p.status & CastleRightsMask
}
func (p *Position) EnPassant() uint64 {
	return p.bitboards[BitEnPassant]
}
func (p *Position) SetTurn(turn uint8) {
	p.status &= ^PlayerMask
	p.status |= turn
}
func (p *Position) SetCastleRights(castleRights uint8) {
	p.status &= ^CastleRightsMask
	p.status |= castleRights
}
func (p *Position) SetEnPassant(rank, file uint8) {
	bit := RFtoB(rank, file)
	p.bitboards[BitEnPassant] = bit
}
func (p *Position) ClearEnPassant() {
	p.bitboards[BitEnPassant] = 0
}
func (p *Position) SetPiece(pieceType uint8, rank, file uint8) {
	bit := RFtoB(rank, file)
	pb, cb := PTtoBB(pieceType)
	p.bitboards[pb] |= bit
	p.bitboards[cb] |= bit
}
func (p *Position) RemovePiece(pieceType uint8, rank, file uint8) {
	notBit := ^RFtoB(rank, file)
	pb, cb := PTtoBB(pieceType)
	p.bitboards[pb] &= notBit
	p.bitboards[cb] &= notBit
}
func (p *Position) MovePiece(fromIndex, toIndex, pieceType uint8) (string, bool) {
	p.ClearEnPassant()
	move, moved := p.specialMove(fromIndex, toIndex, pieceType)
	if !moved {
		fromRank, fromFile := ItoRF(fromIndex)
		toRank, toFile := ItoRF(toIndex)

		// Find piece at destination location
		bit := ItoB(toIndex)
		targetPiece, found := p.identifyPiece(bit)

		// Cannot move on top of own piece
		if found && targetPiece&PlayerMask == pieceType&PlayerMask {
			return "Square is occupied by players own piece", false
		}

		// Determine if piece is a pawn
		piece := pieceType & PieceMask >> 1
		if found && piece == PiecePawn {
			move += FtoN(fromFile)
		} else if piece != PiecePawn { // TODO: sames pieces, same ToIndex - Add N<file>xx, same file, use rank, not file
			move += string(algebraic[piece])
		} else if fromIndex-toIndex == 16 && fromIndex >= 48 && fromIndex <= 55 {
			p.SetEnPassant(fromRank+1, fromFile)
		} else if toIndex-fromIndex == 16 && fromIndex >= 8 && fromIndex <= 15 {
			p.SetEnPassant(fromRank-1, fromFile)
		}

		// Remove piece from source location
		p.RemovePiece(pieceType, fromRank, fromFile)

		// Remove captured piece if found
		if found {
			move += "x"
			p.RemovePiece(targetPiece, toRank, toFile)
		}
		// Set piece at destination location
		p.SetPiece(pieceType, toRank, toFile)
		move += RFtoN(toRank, toFile)
	}

	// Update half move count

	// Update full move count
	if p.Turn() == PlayerBlack {
		p.fullMoves++
	}
	p.SetTurn(1 - p.Turn()) // Switch turn

	// Return move notification
	return move, true
}

func (p *Position) specialMove(fromIndex uint8, toIndex uint8, pieceType uint8) (string, bool) {
	player := pieceType & PlayerMask
	castleRights := p.CastleRights()
	if fromIndex == 60 && toIndex == 62 && player == PlayerWhite {
		if castleRights&CastleRightsWhiteKing != 0 {
			p.SetCastleRights(castleRights & CastleRightsBlackMask)
			p.castle(60, 62, 63, 61, PlayerWhite)
			return "O-O", true
		}
	} else if fromIndex == 60 && toIndex == 58 && player == PlayerWhite {
		if castleRights&CastleRightsWhiteQueen != 0 {
			p.SetCastleRights(castleRights & CastleRightsBlackMask)
			p.castle(60, 58, 56, 59, PlayerWhite)
			return "O-O-O", true
		}
	} else if fromIndex == 4 && toIndex == 6 && player == PlayerBlack {
		if castleRights&CastleRightsBlackKing != 0 {
			p.SetCastleRights(castleRights & CastleRightsWhiteMask)
			p.castle(4, 6, 7, 5, PlayerBlack)
			return "O-O", true
		}
	} else if fromIndex == 4 && toIndex == 2 && player == PlayerBlack {
		if castleRights&CastleRightsBlackQueen != 0 {
			p.SetCastleRights(castleRights & CastleRightsWhiteMask)
			p.castle(4, 2, 0, 3, PlayerBlack)
			return "O-O-O", true
		}
	}
	return "", false
}

func (p *Position) castle(kingFrom, kingTo, rookFrom, rookTo, player uint8) {
	// Move the king
	rank, file := ItoRF(kingTo)
	p.SetPiece(PieceKing|player, rank, file)

	// Move the rook
	rank, file = ItoRF(rookTo)
	p.SetPiece(PieceRook|player, rank, file)

	// Clear the squares
	rank, file = ItoRF(kingFrom)
	p.RemovePiece(PieceKing|player, rank, file)
	rank, file = ItoRF(rookFrom)
	p.RemovePiece(PieceRook|player, rank, file)
}

func (p *Position) ClearSquare(rank, file uint8) {
	notBit := ^RFtoB(rank, file)
	for bb := BitWhite; bb <= BitKings; bb++ {
		p.bitboards[bb] &= notBit
	}
}
func (p *Position) identifyPiece(bit uint64) (uint8, bool) {
	for bb := BitPawns; bb <= BitKings; bb++ {
		if p.bitboards[bb]&bit != 0 {
			pieceType := (bb - 2) << 1
			if p.bitboards[BitBlack]&bit != 0 {
				pieceType |= PlayerBlack
			}
			return pieceType, true
		}
	}
	return 0, false
}

func (p *Position) SetupBoard(fen string) {
	p.fullMoves = 0
	p.halfMoves = make([]uint64, 0)
	p.bitboards = [9]uint64{}

	// Parse the FEN string and set up the board
	index := uint8(0)
bitboards:
	for i := 0; i < len(fen); i++ {
		c := fen[i]
		switch c {
		case '1', '2', '3', '4', '5', '6', '7', '8':
			index += c - '0'
		case '/':
			continue
		case ' ':
			fen = fen[i+1:]
			break bitboards
		case 'K', 'Q', 'R', 'B', 'N', 'P', 'k', 'q', 'r', 'b', 'n', 'p':
			rank, file := ItoRF(index)
			p.SetPiece(fenPieceMap[c], rank, file)
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
func (p *Position) readFenTurn(turn string) {
	t := PlayerWhite
	if strings.EqualFold(turn, "b") {
		t = PlayerBlack
	} else if turn != "w" {
		fmt.Printf("Invalid turn fen: %s\n", turn)
	}
	p.SetTurn(t)
}
func (p *Position) readFenCastling(castling string) {
	castleRights := uint8(0)
	for _, c := range castling {
		switch c {
		case 'K':
			castleRights |= CastleRightsWhiteKing
		case 'Q':
			castleRights |= CastleRightsWhiteQueen
		case 'k':
			castleRights |= CastleRightsBlackKing
		case 'q':
			castleRights |= CastleRightsBlackQueen
		case '-':
			castleRights = 0
			break
		default:
			fmt.Printf("Invalid castling fen: %s\n", castling)
		}
	}
	p.SetCastleRights(castleRights)
}
func (p *Position) readFenEnpassant(enpassant string) {
	switch enpassant {
	case "-":
		p.ClearEnPassant()
	default:
		if rank, file, ok := NtoRF(enpassant); ok {
			p.SetEnPassant(rank, file)
		} else {
			fmt.Printf("Invalid enPassant fen: %s\n", enpassant)
		}
	}
}
func (p *Position) readFenHalfMove(halfMove string) {
	halfMoveCount := 0
	if _, err := fmt.Sscanf(halfMove, "%d", &halfMoveCount); err != nil {
		fmt.Printf("Invalid halfmove fen: %s\n", halfMove)
	}
	p.halfMoves = make([]uint64, halfMoveCount)
}
func (p *Position) readFenFullMove(fullMove string) {
	fullMoveCount := 0
	if _, err := fmt.Sscanf(fullMove, "%d", &fullMoveCount); err != nil {
		fmt.Printf("Invalid fullmove fen: %s\n", fullMove)
	}
	p.fullMoves = fullMoveCount
}

func (p *Position) GenerateFen() string {
	count := 0
	sb := &strings.Builder{}
	for rank := uint8(1); rank <= 8; rank++ {
		for file := uint8(1); file <= 8; file++ {
			bit := RFtoB(rank, file)
			if pieceType, ok := p.identifyPiece(bit); ok {
				p.writeFenEntry(sb, &count, &pieceType)
			} else {
				count++
			}
		}
		if count != 0 {
			p.writeFenEntry(sb, &count, nil)
		}
		if rank < 8 {
			sb.WriteByte('/')
		}
	}
	p.writeFenTurn(sb)
	p.writeFenCastling(sb)
	p.writeFenEnpassant(sb)
	p.writeFenHalfMove(sb)
	p.writeFenFullMove(sb)
	return sb.String()
}
func (p *Position) writeFenEntry(sb *strings.Builder, count *int, pieceType *uint8) {
	if *count > 0 {
		sb.WriteByte(byte(*count + '0'))
	}
	if pieceType != nil {
		sb.WriteByte(fen[*pieceType])
	}
	*count = 0
}
func (p *Position) writeFenTurn(sb *strings.Builder) {
	if p.Turn() == PlayerWhite {
		sb.WriteString(" w")
	} else {
		sb.WriteString(" b")
	}
}
func (p *Position) writeFenCastling(sb *strings.Builder) {
	sb.WriteByte(' ')
	castlingRights := p.CastleRights()
	if castlingRights == 0 {
		sb.WriteByte('-')
		return
	}
	if castlingRights&CastleRightsWhiteKing > 0 {
		sb.WriteByte('K')
	}
	if castlingRights&CastleRightsWhiteQueen > 0 {
		sb.WriteByte('Q')
	}
	if castlingRights&CastleRightsBlackKing > 0 {
		sb.WriteByte('k')
	}
	if castlingRights&CastleRightsBlackQueen > 0 {
		sb.WriteByte('q')
	}
}
func (p *Position) writeFenEnpassant(sb *strings.Builder) {
	ep := p.EnPassant()
	if ep == 0 {
		sb.WriteString(" -")
	} else {
		n := RFtoN(BtoRF(p.EnPassant()))
		sb.WriteString(" " + n)
	}
}
func (p *Position) writeFenHalfMove(sb *strings.Builder) {
	sb.WriteByte(' ')
	sb.WriteString(fmt.Sprintf("%d", len(p.halfMoves)))
}
func (p *Position) writeFenFullMove(sb *strings.Builder) {
	sb.WriteByte(' ')
	sb.WriteString(fmt.Sprintf("%d", p.fullMoves))
}

func (p *Position) debugPrintBoard() {
	for rank := uint8(8); rank >= 1; rank-- {
		fmt.Printf("%d  ", rank)
		for file := uint8(1); file <= 8; file++ {
			str := ". "
			bit := RFtoB(rank, file)
			for bb := BitPawns; bb <= BitKings; bb++ {
				if p.bitboards[bb]&bit != 0 {
					pieceType := (bb - 2) << 1
					if p.bitboards[BitBlack]&bit != 0 {
						pieceType |= PlayerBlack
					}
					str = fmt.Sprintf("%s ", string([]byte{fen[pieceType]}))
					break
				}
			}
			fmt.Print(str)
		}
		fmt.Println()
	}
	fmt.Print("\n   a b c d e f g h\n")
}
func (p *Position) debugPrintBitBoard(bitboard uint64, b uint64) {
	fmt.Printf("Bit Index: %d\n", b)
	for rank := uint8(8); rank >= 1; rank-- {
		fmt.Printf("%d  ", rank)
		for file := uint8(1); file <= 8; file++ {
			bit := RFtoB(rank, file)
			str := ". "
			if bit == b {
				str = "X "
			}
			if bitboard&bit != 0 {
				str = "1 "
				if bit == b {
					str = "@ "
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
	//bs := make([]byte, 8)
	//for i := range keys {
	//	n, err := rand.Read(bs)
	//	if err != nil || n != 8 {
	//		panic("Random key generation failed")
	//	}
	//	keys[i] = binary.BigEndian.Uint64(bs)
	//}
	generateRanksAndFiles()
	generateKnightMoves()
	generateKingMoves()
	generatePawnMoves()
}

func generateRanksAndFiles() {
	for i := uint8(0); i < 8; i++ {
		ranks[i] = 0xFF << (i * 8)
		files[i] = 0x0101010101010101 << i
	}
}

func generateKnightMoves() {
	for rank := uint8(1); rank <= 8; rank++ {
		for file := uint8(1); file <= 8; file++ {
			bit := 63 - RFtoI(rank, file)
			if rank < 7 && file > 1 {
				knightMoves[bit] |= 1 << (bit + 17) //nnw
			}
			if rank < 7 && file < 8 {
				knightMoves[bit] |= 1 << (bit + 15) //nne
			}
			if rank < 8 && file > 2 {
				knightMoves[bit] |= 1 << (bit + 10) //wnw
			}
			if rank < 8 && file < 7 {
				knightMoves[bit] |= 1 << (bit + 6) //ene
			}
			if rank > 1 && file > 2 {
				knightMoves[bit] |= 1 << (bit - 6) //ese
			}
			if rank > 1 && file < 7 {
				knightMoves[bit] |= 1 << (bit - 10) // wsw
			}
			if rank > 2 && file > 1 {
				knightMoves[bit] |= 1 << (bit - 15) // sse
			}
			if rank > 2 && file < 8 {
				knightMoves[bit] |= 1 << (bit - 17) // ssw
			}
			//debugPrintBitBoard(knightMoves[bit], bit)
			//fmt.Println("********************************************")
		}
	}
}
func generateKingMoves() {
	for rank := uint8(1); rank <= 8; rank++ {
		for file := uint8(1); file <= 8; file++ {
			bit := 63 - RFtoI(rank, file)
			if rank < 8 && file > 1 {
				kingMoves[bit] |= 1 << (bit + 9) //nnw
			}
			if rank < 8 {
				kingMoves[bit] |= 1 << (bit + 8) //nnw
			}
			if rank < 8 && file < 8 {
				kingMoves[bit] |= 1 << (bit + 7) //nnw
			}

			if file > 1 {
				kingMoves[bit] |= 1 << (bit + 1) //nnw
			}
			if file < 8 {
				kingMoves[bit] |= 1 << (bit - 1) //nnw
			}

			if rank > 1 && file > 1 {
				kingMoves[bit] |= 1 << (bit - 7) //nnw
			}
			if rank > 1 {
				kingMoves[bit] |= 1 << (bit - 8) //nnw
			}
			if rank > 1 && file < 8 {
				kingMoves[bit] |= 1 << (bit - 9) //nnw
			}
			//debugPrintBitBoard(kingMoves[bit], bit)
			//fmt.Println("********************************************")
		}
	}
}
func generatePawnMoves() {
	for rank := uint8(2); rank <= 7; rank++ {
		for file := uint8(1); file <= 8; file++ {
			bit := 63 - RFtoI(rank, file)
			if file > 1 {
				whitePawnMoves[bit] |= 1 << (bit + 9) //nnw
			}
			if file < 8 {
				whitePawnMoves[bit] |= 1 << (bit + 7) //nnw
			}
			if file > 1 {
				blackPawnMoves[bit] |= 1 << (bit - 7) //nnw
			}
			if file < 8 {
				blackPawnMoves[bit] |= 1 << (bit - 9) //nnw
			}
			//debugPrintBitBoard(whitePawnMoves[bit], uint64(1<<bit))
			//fmt.Println("********************************************")
		}
	}
}

//func(p *Position) Hash() uint64 {
//	var hash uint64
//	for i := range p.bitboards {
//		board := p.bitboards[i]
//		index := 0
//		for board != 0 {
//			if board&1 != 0 {
//				hash ^= keys[index]
//			}
//			board >>= 1
//			index++
//		}
//	}
//	return hash
//}
