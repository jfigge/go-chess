package board

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"strings"
	"us.figge.chess/internal/player"
	. "us.figge.chess/internal/shared"
)

type BoardOptions func(b *Board)

func OptSetup(fen string) BoardOptions {
	return func(b *Board) {
		b.fen = fen
	}
}

type Board struct {
	c         Configuration
	op        *ebiten.DrawImageOptions
	img       *ebiten.Image
	players   [2]*player.Player
	turn      uint8
	board     [64]uint8
	enpassant uint8
	fullMove  uint
	halfMove  uint8
	fen       string
}

func NewBoard(c Configuration, options ...BoardOptions) *Board {
	op := &ebiten.DrawImageOptions{}
	board := &Board{
		c:  c,
		op: op,
		players: [2]*player.Player{
			player.NewPlayer(c, White),
			player.NewPlayer(c, Black),
		},
		enpassant: 0xff,
		turn:      White,
		halfMove:  0,
		fen:       "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
	}
	for _, option := range options {
		option(board)
	}
	board.makeBoardImage()
	board.SetFen(board.fen)
	return board
}

func (b *Board) Draw(target *ebiten.Image) {
	// Draw the board
	img := ebiten.NewImageFromImage(b.img)
	b.players[0].Draw(img)
	b.players[1].Draw(img)
	target.DrawImage(img, b.op)
}

func (b *Board) makeBoardImage() {
	squareSize := int(b.c.SquareSize())
	white := b.c.WhiteColor()
	black := b.c.BlackColor()
	b.img = ebiten.NewImage(squareSize*8, squareSize*8)
	vector.DrawFilledRect(b.img, 0, 0, float32(squareSize*8), float32(squareSize*8), black, false)
	s := float32(squareSize)
	for i := 0; i < 8; i += 2 {
		for j := 0; j < 8; j += 2 {
			vector.DrawFilledRect(b.img, float32(i*squareSize), float32(j*squareSize), s, s, white, false)
			vector.DrawFilledRect(b.img, float32((i+1)*squareSize), float32((j+1)*squareSize), s, s, white, false)
		}
	}
}

func (b *Board) Fen() string {
	return ""
}

func (b *Board) SetFen(fen string) {
	b.fen = fen
	b.board = [64]uint8{}
	index := 0
	b.turn = White
	b.players[0] = player.NewPlayer(b.c, White)
	b.players[1] = player.NewPlayer(b.c, Black)
	for index < 64 && fen != "" {
		p := fen[0]
		fen = fen[1:]
		switch p {
		case '1', '2', '3', '4', '5', '6', '7', '8':
			squares := int(p - '0')
			for i := 0; i < squares; i++ {
				b.board[index] = 0
				index++
			}
		case '/':
			continue
		case ' ':
			break
		case 'K', 'Q', 'R', 'B', 'N', 'P':
			b.board[index] = FenPieceMap[p]
			b.players[0].AddPiece(FenPieceMap[p], uint8(index%8+1), uint8((7-index/8)+1))
			index++
		case 'k', 'q', 'r', 'b', 'n', 'p':
			b.board[index] = FenPieceMap[p]
			b.players[1].AddPiece(FenPieceMap[p], uint8(index%8+1), uint8((7-index/8)+1))
			index++
		}
	}
	parts := strings.Split(fen, " ")
	if len(parts) > 0 {
		b.setTurn(parts[0])
		parts = parts[1:]
	}
	if len(parts) > 0 {
		b.setCastling(parts[0], b.players[0], b.players[1])
		parts = parts[1:]
	}
	if len(parts) > 0 {
		b.setEnpassant(parts[0])
		parts = parts[1:]
	}
	if len(parts) > 0 {
		b.setHalfMove(parts[0])
		parts = parts[1:]
	}
	if len(parts) > 0 {
		b.setFullMove(parts[0])
		parts = parts[1:]
	}

	fmt.Printf("Remaining fen: %s\n", fen)
}

func (b *Board) setTurn(turn string) {
	b.turn = White
	if turn == "b" {
		b.turn = Black
	} else if turn != "w" {
		fmt.Printf("Invalid turn fen: %s", turn)
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
			fmt.Printf("Invalid castling fen: %s", castling)
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
		b.enpassant = uint8((rank-1)*8 + file)
	}
}
func (b *Board) setHalfMove(halfMove string) {
	halfMoveCount := uint8(0)
	if _, err := fmt.Sscanf(halfMove, "%d", &halfMoveCount); err != nil {
		fmt.Printf("Invalid halfmove fen: %s", halfMove)
	}
	b.halfMove = halfMoveCount
}
func (b *Board) setFullMove(fullMove string) {
	fullMoveCount := uint(0)
	if _, err := fmt.Sscanf(fullMove, "%d", &fullMoveCount); err != nil {
		fmt.Printf("Invalid fullmove fen: %s", fullMove)
	}
	b.fullMove = fullMoveCount
}
