package graphics

import (
	"embed"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	. "us.figge.chess/internal/common"
)

type Piece struct {
	img       *ebiten.Image
	pieceType uint8
	name      string
	color     string
}

var (
	names    = [6]string{"Pawn", "Knight", "Bishop", "Rook", "Queen", "King"}
	colors   = [2]string{"White", "Black"}
	imageMap = map[uint8]uint8{
		PiecePawn:   0b01010000,
		PieceKnight: 0b00110000,
		PieceBishop: 0b00100000,
		PieceRook:   0b01000000,
		PieceQueen:  0b00010000,
		PieceKing:   0b00000000,
	}
	pieces map[uint8]*Piece
)

var (
	//go:embed assets/*.png
	assets embed.FS
	sheet  = mustLoadImage("assets/pieces2.png")
)

func mustLoadImage(name string) *ebiten.Image {
	f, err := assets.Open(name)
	if err != nil {
		panic(err)
	}
	defer func() { _ = f.Close() }()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	return ebiten.NewImageFromImage(img)
}

func InitPieces(squareSize int) {
	sheetSize := sheet.Bounds().Size()
	scaleX := float64(squareSize) / float64(sheetSize.X/6)
	scaleY := float64(squareSize) / float64(sheetSize.Y/2)
	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
	op.GeoM.Scale(scaleX, scaleY)
	pieces = make(map[uint8]*Piece)

	for pieceType, location := range imageMap {
		for colorType := PlayerWhite; colorType <= PlayerBlack; colorType++ {
			pieces[pieceType|colorType] = makeEntity(op, squareSize, sheetSize, location|pieceType|colorType)
		}
	}
}

func makeEntity(op *ebiten.DrawImageOptions, squareSize int, sheetSize image.Point, data uint8) *Piece {
	p := &Piece{
		img:       ebiten.NewImage(squareSize, squareSize),
		name:      names[(data&PieceMask)>>1],
		color:     colors[data&PlayerMask],
		pieceType: data & (PieceMask | PlayerMask),
	}
	x := int(data >> 4)
	y := int(data & PlayerMask)
	sx := sheetSize.X / 6
	sy := sheetSize.Y / 2
	p.img.DrawImage(sheet.SubImage(image.Rect(x*sx, y*sy, (x+1)*sx, (y+1)*sy)).(*ebiten.Image), op)
	return p
}

func GetPiece(pieceType uint8) *Piece {
	return pieces[pieceType]
}

func (p *Piece) Draw(dst *ebiten.Image, op *ebiten.DrawImageOptions) {
	dst.DrawImage(p.img, op)
}

func (p *Piece) Name() string      { return p.name }
func (p *Piece) ColorName() string { return p.color }
func (p *Piece) Type() uint8       { return p.pieceType }
