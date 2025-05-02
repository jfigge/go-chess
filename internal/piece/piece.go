package piece

import (
	"github.com/hajimehoshi/ebiten/v2"
	. "us.figge.chess/internal/shared"
)

type Piece struct {
	Configuration
	Token
	rank     int
	file     int
	dragging bool
	dx       float64
	dy       float64
	op       *ebiten.DrawImageOptions
}

func NewPiece(c Configuration, pieceType uint8) *Piece {
	piece := &Piece{
		Configuration: c,
		Token:         c.Token(pieceType),
		rank:          0,
		file:          0,
		op:            &ebiten.DrawImageOptions{},
	}
	return piece
}

func (p *Piece) Draw(target *ebiten.Image) {
	if p.dragging {
		p.op.GeoM.Reset()
		x, y := ebiten.CursorPosition()
		p.op.GeoM.Translate(float64(x)-p.dx, float64(y)-p.dy)
	}
	p.Token.Draw(target, p.op)
}

func (p *Piece) Rank() int {
	return p.rank
}
func (p *Piece) File() int {
	return p.file
}

func (p *Piece) Position(rank, file int) {
	p.rank = rank
	p.file = file
	p.op.GeoM.Reset()
	x, y := p.TranslateRFtoXY(rank, file)
	p.op.GeoM.Translate(x, y)
	p.dragging = false
}

func (p *Piece) StartDrag() {
	if !p.dragging {
		p.dragging = true
		x, y := ebiten.CursorPosition()
		originX, originY := p.TranslateRFtoXY(p.rank, p.file)
		p.dx = float64(x) - originX
		p.dy = float64(y) - originY
	}
}

func (p *Piece) IsDragging() bool {
	return p.dragging
}

func (p *Piece) StopDrag() {
	if p.dragging {
		p.dragging = true
		p.Position(p.rank, p.file)
	}
}
