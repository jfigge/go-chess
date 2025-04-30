package board

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"us.figge.chess/internal/player"
	"us.figge.chess/internal/shared"
)

type Board struct {
	c       shared.Configuration
	img     *ebiten.Image
	players [2]*player.Player
	op      *ebiten.DrawImageOptions
}

func NewBoard(c shared.Configuration) *Board {
	op := &ebiten.DrawImageOptions{}
	board := &Board{
		c:  c,
		op: op,
		players: [2]*player.Player{
			player.NewPlayer(c, player.White),
			player.NewPlayer(c, player.Black),
		},
	}
	board.makeBoardImage()
	return board
}

func (b *Board) Draw(screen *ebiten.Image) {
	// Draw the board
	img := ebiten.NewImageFromImage(b.img)
	b.players[0].Draw(img)
	b.players[1].Draw(img)
	screen.DrawImage(img, b.op)
}

func (b *Board) makeBoardImage() {
	squareSize := b.c.SquareSize()
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
