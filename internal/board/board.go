package board

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"strconv"
	"us.figge.chess/internal/piece"
	"us.figge.chess/internal/player"
	. "us.figge.chess/internal/shared"
)

type square struct {
	piece      *piece.Piece
	background color.Color
	highlight  bool
	valid      bool
	invalid    bool
	x          float32
	y          float32
	size       float32
}

type Board struct {
	Configuration
	players        [2]*player.Player
	squares        [64]square
	turn           uint8
	enpassant      int
	fullMove       int
	halfMove       int
	fen            string
	mouseDown      bool
	mouseFirstDown bool
	dragPiece      *piece.Piece
	dragIndex      int
}

func NewBoard(c Configuration, options ...BoardOptions) *Board {
	board := &Board{
		Configuration: c,
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
	board.SetFen(board.fen)
	return board
}

func (b *Board) Update() {
	b.mouseFirstDown = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	b.mouseDown = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	//b.mouseFirstDown = false
	//down := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	//if down && !b.mouseDown {
	//	fmt.Println("Left mouse button pressed")
	//	b.mouseFirstDown = true
	//	b.mouseDown = true
	//} else if !down && b.mouseDown {
	//	fmt.Println("Left mouse button released")
	//	b.mouseDown = false
	//}
}

func (b *Board) Draw(target *ebiten.Image) {
	var ok bool
	var rank int
	var file int
	cursor := -1
	x, y := ebiten.CursorPosition()
	if rank, file, ok = b.TranslateXYtoRF(x, y); ok {
		cursor = b.TranslateRFtoIndex(rank, file)
	}
	for i := 0; i < 64; i++ {
		highlighted := cursor == i
		p := b.squares[i].piece
		if highlighted && b.mouseFirstDown && p != nil && b.dragPiece == nil {
			b.dragPiece = p
			b.dragIndex = i
			p.StartDrag()
		} else if b.dragPiece != nil && !b.mouseDown && (cursor == i || cursor == 0xff) {
			if p == nil {
				p = b.dragPiece
				b.squares[i].piece = b.dragPiece
				b.squares[b.dragIndex].piece = nil
				p.Position(rank, file)
				b.fen = b.Fen()
			}
			b.dragPiece.StopDrag()
			b.dragPiece = nil
			b.dragIndex = -1
		}
		b.squares[i].Draw(b, target, highlighted)
		if p != nil && i != b.dragIndex {
			p.Draw(target)
		}
	}
	if b.dragPiece != nil {
		b.dragPiece.Draw(target)
	}
	if b.EnableDebug() {
		if cursor != -1 {
			ebitenutil.DebugPrintAt(target, "Rank: "+b.TranslateRFtoN(rank, file), b.DebugX(0), b.DebugY())
			ebitenutil.DebugPrintAt(target, "Index: "+strconv.Itoa(int(cursor)), b.DebugX(1), b.DebugY())
			p := b.squares[cursor].piece
			if p != nil {
				ebitenutil.DebugPrintAt(target, p.Color(), b.DebugX(2), b.DebugY())
				ebitenutil.DebugPrintAt(target, p.Name(), b.DebugX(3), b.DebugY())
			}
		}
		turn := "White"
		if b.turn == Black {
			turn = "Black"
		}
		ebitenutil.DebugPrintAt(target, "Turn: "+turn, b.DebugX(4), b.DebugY())
		ebitenutil.DebugPrintAt(target, "Fen: "+b.fen, b.DebugX(0), b.DebugFen())
	}
}

func (s *square) Draw(b *Board, target *ebiten.Image, highlight bool) {
	c := s.background
	if highlight {
		c = b.ColorHighlight()
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			c = b.ColorInvalid()
		}
	} else if s.valid {
		c = b.ColorValid()
	} else if s.invalid {
		c = b.ColorInvalid()
	}
	vector.DrawFilledRect(target, s.x, s.y, s.size, s.size, c, false)
}

func (b *Board) resetBoard() {
	b.players[0] = player.NewPlayer(Configuration(b), White)
	b.players[1] = player.NewPlayer(Configuration(b), Black)
	b.turn = White
	b.fullMove = 0
	b.halfMove = 0
	b.enpassant = 0xff
	b.dragIndex = 0xff
	b.dragPiece = nil
	for i := 0; i < 64; i++ {
		b.squares[i] = square{
			background: b.ColorBlack(),
			size:       float32(b.SquareSize()),
		}
		if i%2 == (i/8)%2 {
			b.squares[i].background = b.ColorWhite()
		}
		x, y := b.TranslateIndexToXY(i)
		b.squares[i].x = float32(x)
		b.squares[i].y = float32(y)
	}
}
